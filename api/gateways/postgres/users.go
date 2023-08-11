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
		dbUser.EnsAvatarOriginalUrl,
		dbUser.EnsAvatarOriginalContentType,
		dbUser.EnsAvatarFormattedUrl,
		dbUser.EnsAvatarFormattedContentType,
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
	originalURL := pgtype.Text{}
	originalContentType := pgtype.Text{}
	formattedURL := pgtype.Text{}
	formattedContentType := pgtype.Text{}
	if avatar != nil {
		fileName.String = avatar.FileName()
		fileName.Valid = true
		originalURL.String = avatar.OriginalURL()
		originalURL.Valid = true
		originalContentType.String = avatar.OriginalContentType()
		originalContentType.Valid = true
		formattedURL.String = avatar.FormattedURL()
		formattedURL.Valid = true
		formattedContentType.String = avatar.FormattedContentType()
		formattedContentType.Valid = true
	} else {
		fileName.Valid = false
		originalURL.Valid = false
		originalContentType.Valid = false
		formattedURL.Valid = false
		formattedContentType.Valid = false
	}
	return p.queries.UpdateUser(ctx, bindings.UpdateUserParams{
		Address:                       address,
		EnsName:                       ensName,
		EnsAvatarFileName:             fileName,
		EnsAvatarOriginalUrl:          originalURL,
		EnsAvatarOriginalContentType:  originalContentType,
		EnsAvatarFormattedUrl:         formattedURL,
		EnsAvatarFormattedContentType: formattedContentType,
	})
}

func toUser(
	address string,
	name pgtype.Text,
	avatarFileName pgtype.Text,
	avatarOriginalUrl pgtype.Text,
	avatarOriginalContentType pgtype.Text,
	avatarFormattedUrl pgtype.Text,
	avatarFormattedContentType pgtype.Text,
	reputation int64,
	createdAt pgtype.Timestamp,
	updatedAt pgtype.Timestamp,
) entities.User {
	var ensName *string
	if name.Valid {
		ensName = &name.String
	}
	var ensAvatar *entities.Image
	if avatarOriginalUrl.Valid {
		avatar := entities.NewImage(
			avatarFileName.String,
			avatarOriginalUrl.String,
			avatarOriginalContentType.String,
			avatarFormattedUrl.String,
			avatarFormattedContentType.String,
		)
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
