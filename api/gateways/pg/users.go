package pg

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/db/bindings"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p *postgresGateway) CreateOrUpdateUser(ctx context.Context, address string, ensName *string) (entities.User, error) {
	var sqlEnsName pgtype.Text
	if ensName != nil {
		sqlEnsName = pgtype.Text{
			String: *ensName,
			Valid:  true,
		}
	} else {
		sqlEnsName = pgtype.Text{
			String: "",
			Valid:  false,
		}
	}
	params := bindings.CreateOrUpdateUserParams{
		Address: address,
		EnsName: sqlEnsName,
	}
	user, err := p.queries.CreateOrUpdateUser(ctx, params)

	if err != nil {
		return entities.User{}, err
	}

	return entities.NewUser(entities.UserParams{
		Address:   address,
		EnsName:   ensName,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}), nil
}
