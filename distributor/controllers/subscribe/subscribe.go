package subscribe

import (
	"context"
	"fmt"
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
		results, err := s.readMessages(ctx, group, consumer)

		if err != nil {
			s.logger.Error(ctx).Err(err).Msg("error reading messages from streams")
			continue
		}

		for _, result := range results {
			for _, message := range result.Messages {
				s.logger.Info(ctx).Msgf("received message: %v %v %v", result.Stream, message.ID, message.Values)
			}
		}
	}
}

// Reads messages from the streams starting by checking the pending messages that are unacknowledged
// If there are no messages, block for 10 seconds
func (s *subscriber) readMessages(ctx context.Context, group string, consumer string) ([]redis.XStream, error) {
	for _, stream := range []string{common.VoteStream} {
		messages, _, err := s.client.XAutoClaim(ctx, &redis.XAutoClaimArgs{
			Stream:  stream,
			Group:   group,
			Start:   "0-0",
			MinIdle: time.Minute * 15,
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
		Block:    time.Second * 10,
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
