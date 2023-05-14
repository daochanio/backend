package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/db/bindings"
)

func (p *postgresGateway) CreateThread(ctx context.Context, thread entities.Thread) (entities.Thread, error) {
	createdThread, err := p.queries.CreateThread(ctx, bindings.CreateThreadParams{
		Address:          thread.Address(),
		Title:            thread.Title(),
		Content:          thread.Content(),
		ImageFileName:    thread.Image().FileName(),
		ImageUrl:         thread.Image().Url(),
		ImageContentType: thread.Image().ContentType(),
	})

	if err != nil {
		return entities.Thread{}, err
	}

	var deletedAt *time.Time
	if createdThread.DeletedAt.Valid {
		deletedAt = &createdThread.DeletedAt.Time
	}

	image := entities.NewImage(createdThread.ImageFileName, createdThread.ImageUrl, createdThread.ImageContentType)

	return entities.NewThread(entities.ThreadParams{
		Id:        createdThread.ID,
		Address:   createdThread.Address,
		Title:     createdThread.Title,
		Content:   createdThread.Content,
		Image:     image,
		Votes:     0,
		CreatedAt: createdThread.CreatedAt.Time,
		IsDeleted: createdThread.IsDeleted,
		DeletedAt: deletedAt,
	}), nil
}

func (p *postgresGateway) GetThreads(ctx context.Context, limit int64) ([]entities.Thread, error) {
	threads, err := p.queries.GetThreads(ctx, limit)

	if err != nil {
		return nil, err
	}

	threadEnts := []entities.Thread{}
	for _, thread := range threads {
		var deletedAt *time.Time
		if thread.DeletedAt.Valid {
			deletedAt = &thread.DeletedAt.Time
		}

		image := entities.NewImage(thread.ImageFileName, thread.ImageUrl, thread.ImageContentType)
		entitie := entities.NewThread(entities.ThreadParams{
			Id:        thread.ID,
			Address:   thread.Address,
			Title:     thread.Title,
			Content:   thread.Content,
			Image:     image,
			Votes:     thread.Votes,
			CreatedAt: thread.CreatedAt.Time,
			IsDeleted: thread.IsDeleted,
			DeletedAt: deletedAt,
		})
		threadEnts = append(threadEnts, entitie)
	}
	return threadEnts, nil
}

func (p *postgresGateway) GetThreadById(ctx context.Context, id int64) (entities.Thread, error) {
	thread, err := p.queries.GetThread(ctx, id)

	if errors.Is(err, sql.ErrNoRows) {
		return entities.Thread{}, common.ErrNotFound
	}

	if err != nil {
		return entities.Thread{}, err
	}

	var deletedAt *time.Time
	if thread.DeletedAt.Valid {
		deletedAt = &thread.DeletedAt.Time
	}

	image := entities.NewImage(thread.ImageFileName, thread.ImageUrl, thread.ImageContentType)
	entitie := entities.NewThread(entities.ThreadParams{
		Id:        thread.ID,
		Address:   thread.Address,
		Title:     thread.Title,
		Content:   thread.Content,
		Image:     image,
		Votes:     thread.Votes,
		CreatedAt: thread.CreatedAt.Time,
		IsDeleted: thread.IsDeleted,
		DeletedAt: deletedAt,
	})
	return entitie, nil
}

func (p *postgresGateway) DeleteThread(ctx context.Context, id int64) error {
	_, err := p.queries.DeleteThread(ctx, id)

	if errors.Is(err, sql.ErrNoRows) {
		return common.ErrNotFound
	}

	return err
}
