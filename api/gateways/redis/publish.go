package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
	"github.com/redis/go-redis/v9"
)

func (r *redisGateway) PublishVote(ctx context.Context, vote entities.Vote) error {
	voteJson := common.VoteMessage{
		Id:      vote.Id(),
		Address: vote.Address(),
		Type:    vote.Type(),
	}

	message, err := json.Marshal(voteJson)

	if err != nil {
		return fmt.Errorf("error marshalling vote message: %w", err)
	}

	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: common.VoteStream,
		ID:     "*",
		MaxLen: 10000,
		Values: map[string]any{
			"body": message,
		},
	}).Err()
}
