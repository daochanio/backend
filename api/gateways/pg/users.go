package pg

import (
	"context"
	"database/sql"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/db/bindings"
)

func (p *postgresGateway) CreateOrUpdateUser(ctx context.Context, address string, ensName *string) (entities.User, error) {
	var sqlEnsName sql.NullString
	if ensName != nil {
		sqlEnsName = sql.NullString{
			String: *ensName,
			Valid:  true,
		}
	} else {
		sqlEnsName = sql.NullString{
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
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}), nil
}
