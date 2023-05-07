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
	buffer                map[string]common.VoteMessage
	lastFlush             time.Time
}

func NewSubscriber(logger common.Logger, settings settings.Settings, commonSettings common.CommonSettings, aggregateVotesUseCase *usecases.AggregateVotesUseCase) Subscriber {
	opt, err := redis.ParseURL(settings.StreamConnectionString())

	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	buffer := map[string]common.VoteMessage{}
	lastFlush := time.Now()

	return &subscriber{
		logger,
		settings,
		commonSettings,
		client,
		aggregateVotesUseCase,
		buffer,
		lastFlush,
	}
}

// Subscribe to the vote stream and read incoming votes
// Buffer comments/thread ids of votes as having "dirty" vote counts.
// Flush the buffer at a certain length threshold or past a certain number of seconds.
// Buffering is to avoid excessive writes on the same column for hot threads/comments
// TODO: Could make this more efficient by having the vote aggregation usecase batch update instead of one by one
func (s *subscriber) Start(ctx context.Context) error {
	group := s.commonSettings.Appname()
	consumer := s.commonSettings.Hostname()

	_ = s.client.XGroupCreateMkStream(ctx, common.VoteStream, group, "$").Err()

	for {
		if len(s.buffer) > 1000 || (time.Since(s.lastFlush) > time.Second*15 && len(s.buffer) > 0) {
			s.logger.Info(ctx).Msgf("flushing buffer with %v dirty comments/threads", len(s.buffer))
			s.lastFlush = time.Now()
			for key, voteMessage := range s.buffer {
				if err := s.aggregateVotesUseCase.Execute(ctx, usecases.AggregateVotesInput{
					Id:   voteMessage.Id,
					Type: voteMessage.Type,
				}); err != nil {
					s.logger.Error(ctx).Err(err).Msgf("error aggregating votes for %v %v", voteMessage.Type, voteMessage.Id)
				}
				delete(s.buffer, key)
			}
		}

		results, err := s.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{common.VoteStream, ">"},
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
				if voteMessage, err := s.parseMessage(ctx, result.Stream, message); err != nil {
					s.logger.Error(ctx).Err(err).Msgf("error parsing message: %v %v %v", result.Stream, message.ID, message.Values)
				} else {
					key := fmt.Sprintf("%v:%v", voteMessage.Type, voteMessage.Id)
					s.buffer[key] = voteMessage
				}

				if err := s.client.XAck(ctx, result.Stream, group, message.ID).Err(); err != nil {
					s.logger.Error(ctx).Err(err).Msgf("error acknowledging message: %v %v %v", result.Stream, message.ID, message.Values)
				}
			}
		}
	}
}

func (s *subscriber) parseMessage(ctx context.Context, stream string, message redis.XMessage) (common.VoteMessage, error) {
	if stream != common.VoteStream {
		return common.VoteMessage{}, fmt.Errorf("invalid stream: %v", stream)
	}

	body := []byte(message.Values["body"].(string))
	var voteMessage common.VoteMessage
	if err := json.Unmarshal(body, &voteMessage); err != nil {
		return common.VoteMessage{}, fmt.Errorf("error unmarshalling vote message: %v %w", message, err)
	}

	return voteMessage, nil
}
