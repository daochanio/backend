package postgres

import (
	"context"

	"github.com/daochanio/backend/db/bindings"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p *postgresGateway) UpsertUser(ctx context.Context, address string) error {
	return p.queries.UpsertUser(ctx, address)
}

func (p *postgresGateway) UpdateUser(ctx context.Context, address string, ensName *string) error {
	name := pgtype.Text{}
	if ensName != nil {
		name.String = *ensName
		name.Valid = true
	} else {
		name.Valid = false
	}
	return p.queries.UpdateUser(ctx, bindings.UpdateUserParams{
		Address: address,
		EnsName: name,
	})
}
