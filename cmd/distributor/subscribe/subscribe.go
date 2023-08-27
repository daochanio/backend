package subscribe

import (
	"context"
	"fmt"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/usecases"
	"github.com/redis/go-redis/v9"
)

type Subscriber interface {
	Start(ctx context.Context, config SubscriberConfig)
	Shutdown(ctx context.Context)
}

type subscriber struct {
	logger       common.Logger
	settings     common.Settings
	client       *redis.Client
	processVotes *usecases.ProcessVote
}

type SubscriberConfig struct {
	ConnectionString string
	DialTimeout      time.Duration
	MinIdleConns     int
	PoolSize         int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
}

func NewSubscriber(
	logger common.Logger,
	settings common.Settings,
	processVotes *usecases.ProcessVote,
) Subscriber {
	return &subscriber{
		logger:       logger,
		settings:     settings,
		client:       nil,
		processVotes: processVotes,
	}
}

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

	group := s.settings.Appname()
	consumer := s.settings.Hostname()

	_ = s.client.XGroupCreateMkStream(ctx, common.VoteStream, group, "$").Err()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info(ctx).Msg("subscriber stopped")
			return
		default:
			s.execute(ctx, group, consumer)
		}
	}
}

func (s *subscriber) Shutdown(ctx context.Context) {
	s.logger.Info(ctx).Msg("shutting down subscriber")
	if err := s.client.Close(); err != nil {
		s.logger.Error(ctx).Err(err).Msg("error closing subscriber client")
	}
}

func (s *subscriber) execute(ctx context.Context, group string, consumer string) {
	results, err := s.readMessages(ctx, group, consumer)

	if err != nil {
		s.logger.Error(ctx).Err(err).Msg("error reading messages from streams")
		return
	}

	for _, result := range results {
		for _, message := range result.Messages {
			body := []byte(message.Values["body"].(string))
			switch result.Stream {
			case common.VoteStream:
				if err := s.processVoteMessage(ctx, message.ID, body); err != nil {
					s.logger.Error(ctx).Err(err).Msgf("error processing vote message: %v %v %v", result.Stream, message.ID, message.Values)
					continue // skip acking the message if we failed to process it
				}
			default:
				s.logger.Error(ctx).Msgf("inavlid stream %v", result.Stream)
			}

			if err := s.client.XAck(ctx, result.Stream, group, message.ID).Err(); err != nil {
				s.logger.Error(ctx).Err(err).Msgf("error acknowledging message: %v %v %v", result.Stream, message.ID, message.Values)
			}
		}
	}
}

// Reads messages from the streams starting by checking the pending messages that are unacknowledged
// If there are no messages, block for 10 seconds (long polling)
func (s *subscriber) readMessages(ctx context.Context, group string, consumer string) ([]redis.XStream, error) {
	for _, stream := range []string{common.VoteStream} {
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
		Streams:  []string{common.VoteStream, ">"},
		Block:    time.Second * 5,
		Count:    100,
	}).Result()

	if err == redis.Nil {
		return []redis.XStream{}, nil
	}

	if err != nil {
		return []redis.XStream{}, fmt.Errorf("error reading messages from streams: %w", err)
	}

	return results, err
}

func (s *subscriber) processVoteMessage(ctx context.Context, messageID string, data []byte) error {
	vote, err := common.Unmarshal[common.VoteMessage](data)

	if err != nil {
		return fmt.Errorf("error parsing vote message: %v", messageID)
	}

	s.logger.Info(ctx).Msgf("processing vote: %v %v %v", vote.Id, vote.Type, vote.Address)

	if err := s.processVotes.Execute(ctx); err != nil {
		return fmt.Errorf("error processing vote: %v", messageID)
	}

	return nil
}
