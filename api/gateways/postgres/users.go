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
	dbUser, err := p.queries.GetUser(ctx, address)

	if errors.Is(err, pgx.ErrNoRows) {
		return entities.User{}, common.ErrNotFound
	}

	if err != nil {
		return entities.User{}, err
	}

	user := toUser(
		dbUser.Address,
		dbUser.EnsName,
		dbUser.EnsAvatarFileName,
		dbUser.EnsAvatarUrl,
		dbUser.EnsAvatarContentType,
		dbUser.Reputation,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
	)

	return user, nil
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

func toUser(
	address string,
	name pgtype.Text,
	avatarFileName pgtype.Text,
	avatarUrl pgtype.Text,
	avatarContentType pgtype.Text,
	reputation int64,
	createdAt pgtype.Timestamp,
	updatedAt pgtype.Timestamp,
) entities.User {
	var ensName *string
	if name.Valid {
		ensName = &name.String
	}
	var ensAvatar *entities.Image
	if avatarUrl.Valid {
		avatar := entities.NewImage(avatarFileName.String, avatarUrl.String, avatarContentType.String)
		ensAvatar = &avatar
	}
	var updatedAtTime *time.Time
	if updatedAt.Valid {
		updatedAtTime = &updatedAt.Time
	}
	return entities.NewUser(entities.UserParams{
		Address:    address,
		EnsName:    ensName,
		EnsAvatar:  ensAvatar,
		Reputation: reputation,
		CreatedAt:  createdAt.Time,
		UpdatedAt:  updatedAtTime,
	})
}
