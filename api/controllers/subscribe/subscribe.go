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
// TODO: Could make this more efficient by having the vote aggregation usecase batch update instead of one by one
func (s *subscriber) Start(ctx context.Context) error {
	group := s.commonSettings.Appname()
	consumer := s.commonSettings.Hostname()

	_ = s.client.XGroupCreateMkStream(ctx, common.SigninStream, group, "$").Err()
	_ = s.client.XGroupCreateMkStream(ctx, common.VoteStream, group, "$").Err()

	for {
		if len(s.messageBuffer) > 1000 || (time.Since(s.lastFlush) > time.Second*60 && len(s.messageBuffer) > 0) {
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

		results, err := s.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{common.SigninStream, common.VoteStream, ">", ">"},
			Block:    time.Second * 10,
		}).Result()

		if err == redis.Nil {
			continue
		}

		if err != nil {
			s.logger.Error(ctx).Err(err).Msg("error reading messages from streams")
			continue
		}

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
}

func parseMessage[T any](ctx context.Context, message redis.XMessage, result *T) error {
	body := []byte(message.Values["body"].(string))
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("error unmarshalling message: %v %w", message, err)
	}

	return nil
}
