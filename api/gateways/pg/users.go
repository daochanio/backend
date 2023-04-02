package pg

import (
	"context"
	"database/sql"

	"github.com/daochanio/backend/db/bindings"
)

func (p *PostgresGateway) CreateOrUpdateUser(ctx context.Context, address string, ensName *string) error {
	params := bindings.CreateOrUpdateUserParams{
		Address: address,
		EnsName: sql.NullString{
			String: *ensName,
			Valid:  ensName != nil,
		},
	}
	return p.queries.CreateOrUpdateUser(ctx, params)
}
