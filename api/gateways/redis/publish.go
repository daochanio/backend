package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
	"github.com/redis/go-redis/v9"
)

func (r *redisStreamGateway) PublishVote(ctx context.Context, vote entities.Vote) error {
	voteJson := common.VoteMessage{
		Id:        vote.Id(),
		Address:   vote.Address(),
		Type:      vote.Type(),
		Value:     vote.Value(),
		UpdatedAt: vote.UpdatedAt(),
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

func (r *redisStreamGateway) PublishSignin(ctx context.Context, address string) error {
	voteJson := common.SigninMessage{
		Address: address,
	}

	message, err := json.Marshal(voteJson)

	if err != nil {
		return fmt.Errorf("error marshalling signin message: %w", err)
	}

	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: common.SigninStream,
		ID:     "*",
		MaxLen: 10000,
		Values: map[string]any{
			"body": message,
		},
	}).Err()
}
