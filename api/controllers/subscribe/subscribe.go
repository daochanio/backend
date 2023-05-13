package subscribe

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	aggregateVotesUseCase *usecases.AggregateVotesUseCase
	hydrateUserUseCase    *usecases.HydrateUserUseCase
	messageBuffer         map[string]*bufferMessage
	lastFlush             time.Time
}

type bufferMessage struct {
	message any
	stream  common.Stream
}

func NewSubscriber(logger common.Logger, settings settings.Settings, client *redis.Client, commonSettings common.CommonSettings, aggregateVotesUseCase *usecases.AggregateVotesUseCase, hydrateUserUseCase *usecases.HydrateUserUseCase) Subscriber {
	buffer := map[string]*bufferMessage{}
	lastFlush := time.Now()
	return &subscriber{
		logger,
		settings,
		commonSettings,
		client,
		aggregateVotesUseCase,
		hydrateUserUseCase,
		buffer,
		lastFlush,
	}
}

// Subscribe to the vote stream and read incoming votes
// Buffer successive messages with the same key to minimize work
// I.e buffering on votes is to avoid excessive writes on the same column for hot threads/comments
// Flush the buffer at a certain length threshold or past a certain number of seconds.
// TODO: if we receive a sigterm and there are messages in the buffer, flush the buffer before exiting
// TODO: Could make this more efficient by processing buffered messages in batches instead of one by one
func (s *subscriber) Start(ctx context.Context) error {
	group := s.commonSettings.Appname()
	consumer := s.commonSettings.Hostname()

	_ = s.client.XGroupCreateMkStream(ctx, common.SigninStream, group, "$").Err()
	_ = s.client.XGroupCreateMkStream(ctx, common.VoteStream, group, "$").Err()

	for {
		if len(s.messageBuffer) >= 1000 || (time.Since(s.lastFlush) > time.Second*15 && len(s.messageBuffer) > 0) {
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

func (s *subscriber) flushBuffer(ctx context.Context) {
	s.logger.Info(ctx).Msgf("flushing buffer with size %v", len(s.messageBuffer))
	s.lastFlush = time.Now()
	for key, bufferMessage := range s.messageBuffer {
		switch bufferMessage.stream {
		case common.VoteStream:
			{
				voteMessage := bufferMessage.message.(common.VoteMessage)
				s.logger.Info(ctx).Msgf("aggregating votes for %v %v", voteMessage.Type, voteMessage.Id)
				if err := s.aggregateVotesUseCase.Execute(ctx, usecases.AggregateVotesInput{
					Id:   voteMessage.Id,
					Type: voteMessage.Type,
				}); err != nil {
					s.logger.Error(ctx).Err(err).Msgf("error aggregating votes for %v %v", voteMessage.Type, voteMessage.Id)
				}
			}
		case common.SigninStream:
			{
				signinMessage := bufferMessage.message.(common.SigninMessage)
				s.logger.Info(ctx).Msgf("hydrating user info for %v", signinMessage.Address)
				if err := s.hydrateUserUseCase.Execute(ctx, usecases.HydrateUserUseCaseInput{
					Address: signinMessage.Address,
				}); err != nil {
					s.logger.Error(ctx).Err(err).Msgf("error hydrating user %v", signinMessage.Address)
				}
			}
		default:
			{
				s.logger.Error(ctx).Msgf("inavlid stream %v", bufferMessage.stream)
			}
		}
		delete(s.messageBuffer, key)
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
			switch result.Stream {
			case common.SigninStream:
				{
					var signinMessage common.SigninMessage
					if err := parseMessage(ctx, message, &signinMessage); err != nil {
						s.logger.Error(ctx).Err(err).Msgf("error processing message: %v %v %v", result.Stream, message.ID, message.Values)
					}
					s.messageBuffer[signinMessage.Address] = &bufferMessage{
						message: signinMessage,
						stream:  result.Stream,
					}
				}
			case common.VoteStream:
				{
					var voteMessage common.VoteMessage
					if err := parseMessage(ctx, message, &voteMessage); err != nil {
						s.logger.Error(ctx).Err(err).Msgf("error processing message: %v %v %v", result.Stream, message.ID, message.Values)
					}
					s.messageBuffer[fmt.Sprintf("%v:%v", voteMessage.Type, voteMessage.Id)] = &bufferMessage{
						message: voteMessage,
						stream:  result.Stream,
					}
				}
			default:
				s.logger.Error(ctx).Msgf("invalid stream: %v", result.Stream)
			}

			if err := s.client.XAck(ctx, result.Stream, group, message.ID).Err(); err != nil {
				s.logger.Error(ctx).Err(err).Msgf("error acknowledging message: %v %v %v", result.Stream, message.ID, message.Values)
			}
		}
	}
}

func parseMessage[T any](ctx context.Context, message redis.XMessage, result *T) error {
	body := []byte(message.Values["body"].(string))
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("error unmarshalling message: %v %w", message, err)
	}

	return nil
}
