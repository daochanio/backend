package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/daochanio/backend/api/entities"
)

const ChallengeNamespace = "challenge"

func (r *redisGateway) GetChallengeByAddress(ctx context.Context, address string) (entities.Challenge, error) {
	message, err := r.client.Get(ctx, getFullKey(ChallengeNamespace, address)).Result()

	if err != nil {
		return entities.Challenge{}, fmt.Errorf("get challenge %w", err)
	}

	return entities.NewChallenge(address, message), nil
}

func (r *redisGateway) SaveChallenge(ctx context.Context, challenge entities.Challenge) error {
	err := r.client.Set(ctx, getFullKey(ChallengeNamespace, challenge.Address()), challenge.Message(), time.Minute*10).Err()

	if err != nil {
		return fmt.Errorf("save challenge %w", err)
	}

	return nil
}
