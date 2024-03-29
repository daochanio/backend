package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/gateways/postgres/bindings"
	"github.com/jackc/pgx/v5"
)

func (p *postgresGateway) CreateThread(
	ctx context.Context,
	address string,
	title string,
	content string,
	image *entities.Image,
) (entities.Thread, error) {
	id, err := p.queries.CreateThread(ctx, bindings.CreateThreadParams{
		Address:                   address,
		Title:                     title,
		Content:                   content,
		ImageFileName:             image.FileName(),
		ImageOriginalUrl:          image.OriginalURL(),
		ImageOriginalContentType:  image.OriginalContentType(),
		ImageFormattedUrl:         image.FormattedURL(),
		ImageFormattedContentType: image.FormattedContentType(),
	})

	if err != nil {
		return entities.Thread{}, err
	}

	return p.GetThreadById(ctx, id)
}

func (p *postgresGateway) GetThreads(ctx context.Context, limit int64) ([]entities.Thread, error) {
	dbThreads, err := p.queries.GetThreads(ctx, limit)

	if err != nil {
		return nil, err
	}

	threads := []entities.Thread{}
	for _, dbThread := range dbThreads {
		var deletedAt *time.Time
		if dbThread.DeletedAt.Valid {
			deletedAt = &dbThread.DeletedAt.Time
		}

		image := entities.NewImage(
			dbThread.ImageFileName,
			dbThread.ImageOriginalUrl,
			dbThread.ImageOriginalContentType,
			dbThread.ImageFormattedUrl,
			dbThread.ImageFormattedContentType,
		)
		user := toUser(
			dbThread.Address,
			dbThread.EnsName,
			dbThread.EnsAvatarFileName,
			dbThread.EnsAvatarOriginalUrl,
			dbThread.EnsAvatarOriginalContentType,
			dbThread.EnsAvatarFormattedUrl,
			dbThread.EnsAvatarFormattedContentType,
			dbThread.Reputation,
			dbThread.UserCreatedAt,
			dbThread.UserUpdatedAt,
		)
		thread := entities.NewThread(entities.ThreadParams{
			Id:        dbThread.ID,
			Title:     dbThread.Title,
			Content:   dbThread.Content,
			Image:     image,
			User:      user,
			Votes:     dbThread.Votes,
			CreatedAt: dbThread.CreatedAt.Time,
			IsDeleted: dbThread.IsDeleted,
			DeletedAt: deletedAt,
		})
		threads = append(threads, thread)
	}
	return threads, nil
}

func (p *postgresGateway) GetThreadById(ctx context.Context, id int64) (entities.Thread, error) {
	dbThread, err := p.queries.GetThread(ctx, id)

	if errors.Is(err, pgx.ErrNoRows) {
		return entities.Thread{}, common.ErrNotFound
	}

	if err != nil {
		return entities.Thread{}, err
	}

	var deletedAt *time.Time
	if dbThread.DeletedAt.Valid {
		deletedAt = &dbThread.DeletedAt.Time
	}

	image := entities.NewImage(
		dbThread.ImageFileName,
		dbThread.ImageOriginalUrl,
		dbThread.ImageOriginalContentType,
		dbThread.ImageFormattedUrl,
		dbThread.ImageFormattedContentType,
	)
	user := toUser(
		dbThread.Address,
		dbThread.EnsName,
		dbThread.EnsAvatarFileName,
		dbThread.EnsAvatarOriginalUrl,
		dbThread.EnsAvatarOriginalContentType,
		dbThread.EnsAvatarFormattedUrl,
		dbThread.EnsAvatarFormattedContentType,
		dbThread.Reputation,
		dbThread.UserCreatedAt,
		dbThread.UserUpdatedAt,
	)
	thread := entities.NewThread(entities.ThreadParams{
		Id:        dbThread.ID,
		Title:     dbThread.Title,
		Content:   dbThread.Content,
		Image:     image,
		User:      user,
		Votes:     dbThread.Votes,
		CreatedAt: dbThread.CreatedAt.Time,
		IsDeleted: dbThread.IsDeleted,
		DeletedAt: deletedAt,
	})
	return thread, nil
}

func (p *postgresGateway) DeleteThread(ctx context.Context, id int64) error {
	_, err := p.queries.DeleteThread(ctx, id)

	if errors.Is(err, pgx.ErrNoRows) {
		return common.ErrNotFound
	}

	return err
}
