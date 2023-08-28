package subscribe

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/domain/usecases"
	"github.com/redis/go-redis/v9"
)

type Subscriber interface {
	Start(ctx context.Context, config SubscriberConfig)
	Shutdown(ctx context.Context)
}

type subscriber struct {
	logger                common.Logger
	client                *redis.Client
	aggregateVotesUseCase *usecases.AggregateVotes
	hydrateUsersUseCase   *usecases.HydrateUsers
	messageBuffer         *[]bufferMessage
	lastFlush             time.Time
}

type bufferMessage struct {
	message redis.XMessage
	stream  redis.XStream
}

type SubscriberConfig struct {
	Group            string
	Consumer         string
	ConnectionString string
	DialTimeout      time.Duration
	MinIdleConns     int
	PoolSize         int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
}

func NewSubscriber(
	logger common.Logger,
	aggregateVotesUseCase *usecases.AggregateVotes,
	hydrateUsersUseCase *usecases.HydrateUsers,
) Subscriber {
	return &subscriber{
		logger:                logger,
		client:                nil,
		aggregateVotesUseCase: aggregateVotesUseCase,
		hydrateUsersUseCase:   hydrateUsersUseCase,
		messageBuffer:         nil,
		lastFlush:             time.Now(),
	}
}

// Subscribe to the vote stream and read incoming votes
// Buffer messages to provide the opportunity for de-duplication of messages with similar keys.
// I.e buffering on votes is to avoid excessive writes on the same column for hot threads/comments or a bad actor writing the same vote over and over.
// Either scenario would cause a lot of write contention on a single row.
// We flush the buffer at a certain length or past a certain number of seconds.
// We can't assume that the messages we are reading in order, since streams are processed in parallel from multiple distributed processes.
// Example: autoclaiming a vote message from the PEL that Node 1 failed to process but is older than a message that Node 2 is currently processing.
func (s *subscriber) Start(ctx context.Context, config SubscriberConfig) {
	s.logger.Info(ctx).Msg("starting redis subscriber")

	opt, err := redis.ParseURL(config.ConnectionString)

	if err != nil {
		panic(err)
	}

	opt.DialTimeout = config.DialTimeout
	opt.MinIdleConns = config.MinIdleConns
	opt.PoolSize = config.PoolSize
	opt.ReadTimeout = config.ReadTimeout
	opt.WriteTimeout = config.WriteTimeout

	s.client = redis.NewClient(opt)

	s.messageBuffer = &[]bufferMessage{}
	s.lastFlush = time.Now()

	_ = s.client.XGroupCreateMkStream(ctx, common.SigninStream, config.Group, "$").Err()
	_ = s.client.XGroupCreateMkStream(ctx, common.VoteStream, config.Group, "$").Err()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info(ctx).Msg("subscriber stopped")
			return
		default:
			s.execute(ctx, config.Group, config.Consumer)
		}
	}
}

func (s *subscriber) Shutdown(ctx context.Context) {
	s.logger.Info(ctx).Msg("shutting down subscriber")

	// ensure we clear the buffer before shutting down
	s.flushBuffer(ctx)

	if err := s.client.Close(); err != nil {
		s.logger.Error(ctx).Err(err).Msg("error closing redis client")
	}
}

func (s *subscriber) execute(ctx context.Context, group string, consumer string) {
	if len(*s.messageBuffer) >= 1000 || (time.Since(s.lastFlush) > time.Second*15 && len(*s.messageBuffer) > 0) {
		s.flushBuffer(ctx)
	}

	results, err := s.readMessages(ctx, group, consumer)

	if err != nil {
		s.logger.Error(ctx).Err(err).Msg("error reading messages from streams")
		return
	}

	s.bufferMessages(ctx, group, results)
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
		Block:    time.Second * 5,
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
		body := []byte(message.Values["body"].(string))
		switch stream {
		case common.VoteStream:
			{
				voteMessage, err := common.Unmarshal[common.VoteMessage](body)
				if err != nil {
					s.logger.Error(ctx).Err(err).Msgf("error parsing vote message: %v %v %v", stream, message.ID, message.Values)
					continue
				}
				vote := entities.NewVote(voteMessage.Id, voteMessage.Address, voteMessage.Value, voteMessage.Type, voteMessage.UpdatedAt)
				votes = append(votes, vote)
			}
		case common.SigninStream:
			{
				signinMessage, err := common.Unmarshal[common.SigninMessage](body)
				if err != nil {
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

	s.hydrateUsersUseCase.Execute(ctx, usecases.HydrateUsersInput{
		Addresses: addresses,
	})
}
