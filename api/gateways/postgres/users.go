package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/db/bindings"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p *postgresGateway) GetUserByAddress(ctx context.Context, address string) (entities.User, error) {
	user, err := p.queries.GetUser(ctx, address)

	if errors.Is(err, pgx.ErrNoRows) {
		return entities.User{}, common.ErrNotFound
	}

	if err != nil {
		return entities.User{}, err
	}

	var ensName *string
	if user.EnsName.Valid {
		ensName = &user.EnsName.String
	}

	var ensAvatar *entities.Image
	if user.EnsAvatarUrl.Valid {
		image := entities.NewImage(user.EnsAvatarFileName.String, user.EnsAvatarUrl.String, user.EnsAvatarContentType.String)
		ensAvatar = &image
	}

	var updatedAt *time.Time
	if user.UpdatedAt.Valid {
		updatedAt = &user.UpdatedAt.Time
	}

	return entities.NewUser(entities.UserParams{
		Address:    user.Address,
		EnsName:    ensName,
		EnsAvatar:  ensAvatar,
		Reputation: user.Reputation,
		CreatedAt:  user.CreatedAt.Time,
		UpdatedAt:  updatedAt,
	}), nil
}

func (p *postgresGateway) UpsertUser(ctx context.Context, address string) error {
	return p.queries.UpsertUser(ctx, address)
}

func (p *postgresGateway) UpdateUser(ctx context.Context, address string, name *string, avatar *entities.Image) error {
	ensName := pgtype.Text{}
	if name != nil {
		ensName.String = *name
		ensName.Valid = true
	} else {
		ensName.Valid = false
	}
	fileName := pgtype.Text{}
	url := pgtype.Text{}
	contentType := pgtype.Text{}
	if avatar != nil {
		fileName.String = avatar.FileName()
		fileName.Valid = true
		url.String = avatar.Url()
		url.Valid = true
		contentType.String = avatar.ContentType()
		contentType.Valid = true
	} else {
		fileName.Valid = false
		url.Valid = false
		contentType.Valid = false
	}
	return p.queries.UpdateUser(ctx, bindings.UpdateUserParams{
		Address:              address,
		EnsName:              ensName,
		EnsAvatarUrl:         url,
		EnsAvatarFileName:    fileName,
		EnsAvatarContentType: contentType,
	})
}
