package pg

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/db/bindings"
)

func (p *postgresGateway) CreateThread(ctx context.Context, address string, title string, content string, imageFileName string, imageURL string, imageContentType string) (int64, error) {
	return p.queries.CreateThread(ctx, bindings.CreateThreadParams{
		Address:          address,
		Title:            title,
		Content:          content,
		ImageFileName:    imageFileName,
		ImageUrl:         imageURL,
		ImageContentType: imageContentType,
	})
}

func (p *postgresGateway) GetThreads(ctx context.Context, limit int32) ([]entities.Thread, error) {
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
			CreatedAt: thread.CreatedAt,
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
		CreatedAt: thread.CreatedAt,
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

func (p *postgresGateway) UpVoteThread(ctx context.Context, id int64, address string) error {
	return p.queries.CreateThreadUpVote(ctx, bindings.CreateThreadUpVoteParams{
		ThreadID: id,
		Address:  address,
	})
}

func (p *postgresGateway) DownVoteThread(ctx context.Context, id int64, address string) error {
	return p.queries.CreateThreadDownVote(ctx, bindings.CreateThreadDownVoteParams{
		ThreadID: id,
		Address:  address,
	})
}

func (p *postgresGateway) UnVoteThread(ctx context.Context, id int64, address string) error {
	return p.queries.CreateThreadUnVote(ctx, bindings.CreateThreadUnVoteParams{
		ThreadID: id,
		Address:  address,
	})
}
