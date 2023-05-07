package subscribe

import (
	"context"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/distributor/settings"
	"github.com/redis/go-redis/v9"
)

type Subscriber interface {
	Start(ctx context.Context) error
}

type subscriber struct {
	logger         common.Logger
	settings       settings.Settings
	commonSettings common.CommonSettings
	client         *redis.Client
}

func NewSubscriber(logger common.Logger, settings settings.Settings, commonSettings common.CommonSettings) Subscriber {
	opt, err := redis.ParseURL(settings.StreamConnectionString())

	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)

	return &subscriber{
		logger,
		settings,
		commonSettings,
		client,
	}
}

// TODO:
// Create a distribution record that will represent the next distribution to execute
//   - Can have things like the transaction id, the block number, the block hash, etc associated with it
//
// As we process votes from the stream, we can hydrate the full vote record by calling the API
// We can then make decisions on whether the vote should be counted towards a distribution or not and create a vote record for it
//   - I.e if the vote is on a comment/thread that is older than a certain cuttoff
//   - We will always write a record to the table for every vote we process, regardless of whether it is counted or not with some kind of accepted/discarded flag
//
// When the next distribution round runs, the records that are accepted and not associated with a distribution can be processed and then tied to a distribution through FK
func (s *subscriber) Start(ctx context.Context) error {
	group := s.commonSettings.Appname()
	consumer := s.commonSettings.Hostname()

	_ = s.client.XGroupCreateMkStream(ctx, common.VoteStream, group, "$").Err()

	for {
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
				s.logger.Info(ctx).Msgf("received message: %v %v %v", result.Stream, message.ID, message.Values)

				if err := s.client.XAck(ctx, result.Stream, group, message.ID).Err(); err != nil {
					s.logger.Error(ctx).Err(err).Msgf("error acknowledging message: %v %v %v", result.Stream, message.ID, message.Values)
				}
			}
		}
	}
}