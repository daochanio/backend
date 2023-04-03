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

func (p *PostgresGateway) CreateThread(ctx context.Context, address string, content string) (int64, error) {
	return p.queries.CreateThread(ctx, bindings.CreateThreadParams{
		Address: address,
		Content: content,
	})
}

func (p *PostgresGateway) GetThreads(ctx context.Context, offset int64, limit int64) ([]entities.Thread, error) {
	threads, err := p.queries.GetThreads(ctx, bindings.GetThreadsParams{
		Offset: int32(offset),
		Limit:  int32(limit),
	})

	if err != nil {
		return nil, err
	}

	threadEnts := []entities.Thread{}
	for _, thread := range threads {
		var deletedAt *time.Time
		if thread.DeletedAt.Valid {
			deletedAt = &thread.DeletedAt.Time
		}

		entitie := entities.
			NewThread().
			SetId(thread.ID).
			SetAddress(thread.Address).
			SetContent(thread.Content).
			SetVotes(thread.Votes).
			SetCreatedAt(thread.CreatedAt).
			SetDeletedAt(deletedAt).
			SetIsDeleted(thread.IsDeleted)
		threadEnts = append(threadEnts, entitie)
	}
	return threadEnts, nil
}

func (p *PostgresGateway) GetThreadById(ctx context.Context, id int64) (entities.Thread, error) {
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

	entitie := entities.
		NewThread().
		SetId(thread.ID).
		SetAddress(thread.Address).
		SetContent(thread.Content).
		SetVotes(thread.Votes).
		SetCreatedAt(thread.CreatedAt).
		SetDeletedAt(deletedAt).
		SetIsDeleted(thread.IsDeleted)
	return entitie, nil
}

func (p *PostgresGateway) DeleteThread(ctx context.Context, id int64) error {
	_, err := p.queries.DeleteThread(ctx, id)

	if errors.Is(err, sql.ErrNoRows) {
		return common.ErrNotFound
	}

	return err
}

func (p *PostgresGateway) UpVoteThread(ctx context.Context, id int64, address string) error {
	return p.queries.CreateThreadUpVote(ctx, bindings.CreateThreadUpVoteParams{
		ThreadID: id,
		Address:  address,
	})
}

func (p *PostgresGateway) DownVoteThread(ctx context.Context, id int64, address string) error {
	return p.queries.CreateThreadDownVote(ctx, bindings.CreateThreadDownVoteParams{
		ThreadID: id,
		Address:  address,
	})
}

func (p *PostgresGateway) UnVoteThread(ctx context.Context, id int64, address string) error {
	return p.queries.CreateThreadUnVote(ctx, bindings.CreateThreadUnVoteParams{
		ThreadID: id,
		Address:  address,
	})
}
