package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/gateways/postgres/bindings"
)

func (p *postgresGateway) GetChallengeByAddress(ctx context.Context, address string) (entities.Challenge, error) {
	challenge, err := p.queries.GetChallenge(ctx, address)

	if err != nil {
		return entities.Challenge{}, common.ErrNotFound
	}

	if expiresAt := time.Unix(challenge.ExpiresAt, 0); expiresAt.Before(time.Now()) {
		return entities.Challenge{}, fmt.Errorf("challenge expired")
	}

	return entities.NewChallenge(address, challenge.Message), nil
}

func (p *postgresGateway) SaveChallenge(ctx context.Context, challenge entities.Challenge) error {
	err := p.queries.UpdateChallenge(ctx, bindings.UpdateChallengeParams{
		Address:   challenge.Address(),
		Message:   challenge.Message(),
		ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
	})

	if err != nil {
		return fmt.Errorf("save challenge %w", err)
	}

	return nil
}
