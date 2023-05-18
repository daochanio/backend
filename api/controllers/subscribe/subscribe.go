package subscribe

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
	"github.com/redis/go-redis/v9"
)

type Subscriber interface {
	Start(ctx context.Context) error
}

type subscriber struct {
	logger                common.Logger
	settings              settings.Settings
	commonSettings        common.CommonSettings
	client                *redis.Client
	aggregateVotesUseCase *usecases.AggregateVotes
	hydrateUserUseCase    *usecases.HydrateUsers
	messageBuffer         *[]bufferMessage
	lastFlush             time.Time
}

type bufferMessage struct {
	message redis.XMessage
	stream  redis.XStream
}

func NewSubscriber(logger common.Logger, settings settings.Settings, commonSettings common.CommonSettings, aggregateVotesUseCase *usecases.AggregateVotes, hydrateUsersUseCase *usecases.HydrateUsers) Subscriber {
	client := redis.NewClient(settings.GlobalRedisOptions())
	buffer := &[]bufferMessage{}
	lastFlush := time.Now()
	return &subscriber{
		logger,
		settings,
		commonSettings,
		client,
		aggregateVotesUseCase,
		hydrateUsersUseCase,
		buffer,
		lastFlush,
	}
}

// Subscribe to the vote stream and read incoming votes
// Buffer successive messages with the same key to minimize work
// I.e buffering on votes is to avoid excessive writes on the same column for hot threads/comments or a bad actor writing the same vote over and over.
// Either scenario would cause a lot of write contention on a single row.
// We flush the buffer at a certain length or past a certain number of seconds.
// We can't assume that the messages we are reading in order, since streams are processed in parallel from multiple distributed processes.
// Example: autoclaiming a vote message from the PEL that Node 1 failed to process but is older than a message that Node 2 is currently processing.
//
// TODO: If we receive a sigterm and there are messages in the buffer, flush the buffer and ensure it finishes before exiting
func (s *subscriber) Start(ctx context.Context) error {
	group := s.commonSettings.Appname()
	consumer := s.commonSettings.Hostname()

	_ = s.client.XGroupCreateMkStream(ctx, common.SigninStream, group, "$").Err()
	_ = s.client.XGroupCreateMkStream(ctx, common.VoteStream, group, "$").Err()

	for {
		if len(*s.messageBuffer) >= 1000 || (time.Since(s.lastFlush) > time.Second*15 && len(*s.messageBuffer) > 0) {
			s.flushBuffer(ctx)
		}

		results, err := s.readMessages(ctx, group, consumer)

		if err != nil {
			s.logger.Error(ctx).Err(err).Msg("error reading messages from streams")
			continue
		}

		s.bufferMessages(ctx, group, results)
	}
}

// Reads messages from the streams starting by checking the pending messages that are unacknowledged
// If there are no messages, block for 10 seconds
func (s *subscriber) readMessages(ctx context.Context, group string, consumer string) ([]redis.XStream, error) {
	for _, stream := range []string{common.SigninStream, common.VoteStream} {
		messages, _, err := s.client.XAutoClaim(ctx, &redis.XAutoClaimArgs{
			Stream:  stream,
			Group:   group,
			Start:   "0-0",
			MinIdle: time.Minute * 5,
			Count:   1000, // pending entries list has a max size of 1000
		}).Result()

		if err != nil && err != redis.Nil {
			return []redis.XStream{}, fmt.Errorf("error claiming pending messages from stream: %v %w", stream, err)
		}

		if len(messages) > 0 {
			s.logger.Info(ctx).Msgf("claimed %v pending messages from stream %v group %v", len(messages), stream, group)

			return []redis.XStream{{
				Stream:   stream,
				Messages: messages,
			}}, nil
		}
	}

	results, err := s.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{common.SigninStream, common.VoteStream, ">", ">"},
		Block:    time.Second * 10,
		Count:    100,
	}).Result()

	if err == redis.Nil {
		return []redis.XStream{}, nil
	}

	if err != nil {
		return []redis.XStream{}, err
	}

	return results, err
}

func (s *subscriber) bufferMessages(ctx context.Context, group string, results []redis.XStream) {
	for _, result := range results {
		for _, message := range result.Messages {
			messageBuffer := append(*s.messageBuffer, bufferMessage{
				message: message,
				stream:  result,
			})
			s.messageBuffer = &messageBuffer
			if err := s.client.XAck(ctx, result.Stream, group, message.ID).Err(); err != nil {
				s.logger.Error(ctx).Err(err).Msgf("error acknowledging message: %v %v %v", result.Stream, message.ID, message.Values)
			}
		}
	}
}

func (s *subscriber) flushBuffer(ctx context.Context) {
	s.logger.Info(ctx).Msgf("flushing buffer with size %v", len(*s.messageBuffer))
	userAddresses := []string{}
	votes := []entities.Vote{}
	for _, bufferMessage := range *s.messageBuffer {
		stream := bufferMessage.stream.Stream
		message := bufferMessage.message
		switch stream {
		case common.VoteStream:
			{
				var voteMessage common.VoteMessage
				if err := parseMessage(ctx, message, &voteMessage); err != nil {
					s.logger.Error(ctx).Err(err).Msgf("error parsing vote message: %v %v %v", stream, message.ID, message.Values)
					continue
				}
				vote := entities.NewVote(voteMessage.Id, voteMessage.Address, voteMessage.Value, voteMessage.Type, voteMessage.UpdatedAt)
				votes = append(votes, vote)
			}
		case common.SigninStream:
			{
				var signinMessage common.SigninMessage
				if err := parseMessage(ctx, message, &signinMessage); err != nil {
					s.logger.Error(ctx).Err(err).Msgf("error parsing signin message: %v %v %v", stream, message.ID, message.Values)
					continue
				}
				userAddresses = append(userAddresses, signinMessage.Address)
			}
		default:
			{
				s.logger.Error(ctx).Msgf("inavlid stream %v", bufferMessage.stream)
			}
		}
	}

	// wipe buffer/lastFlush regardless of success/failure of below processing to avoid memory leaks from an infinitely growing buffer
	s.lastFlush = time.Now()
	s.messageBuffer = &[]bufferMessage{}

	var wg sync.WaitGroup
	wg.Add(2)

	go s.aggregateVotes(ctx, &wg, votes)
	go s.hydrateUsers(ctx, &wg, userAddresses)

	wg.Wait()
}

func parseMessage[T any](ctx context.Context, message redis.XMessage, result *T) error {
	body := []byte(message.Values["body"].(string))
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("error unmarshalling message: %v %w", message, err)
	}

	return nil
}

func (s *subscriber) aggregateVotes(ctx context.Context, wg *sync.WaitGroup, votes []entities.Vote) {
	defer wg.Done()

	if len(votes) == 0 {
		return
	}

	s.aggregateVotesUseCase.Execute(ctx, usecases.AggregateVotesInput{
		Votes: votes,
	})
}

func (s *subscriber) hydrateUsers(ctx context.Context, wg *sync.WaitGroup, addresses []string) {
	defer wg.Done()

	if len(addresses) == 0 {
		return
	}

	s.hydrateUserUseCase.Execute(ctx, usecases.HydrateUsersInput{
		Addresses: addresses,
	})
}
