package postgres

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/db/bindings"
	"github.com/jackc/pgx/v5/pgtype"
)

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
