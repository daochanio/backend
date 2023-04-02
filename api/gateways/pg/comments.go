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

func (p *PostgresGateway) CreateComment(ctx context.Context, threadId int32, address string, parentCommentId *int32, content string) (int32, error) {
	// begin tx
	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	qtx := p.queries.WithTx(tx)

	id, err := p.queries.CreateComment(ctx, bindings.CreateCommentParams{
		ThreadID: threadId,
		Address:  address,
		Content:  content,
	})

	if err != nil {
		return 0, err
	}

	if err := qtx.CreateSelfClosure(ctx, id); err != nil {
		return 0, err
	}

	// only create parent closures if comment is responding to another comment
	if parentCommentId != nil {
		if err := qtx.CreateParentClosures(ctx, bindings.CreateParentClosuresParams{
			ChildID:  id,
			ParentID: *parentCommentId,
		}); err != nil {
			return 0, err
		}
	}

	return id, tx.Commit()
}

func (p *PostgresGateway) GetComments(ctx context.Context, threadId int32) ([]entities.Comment, error) {
	comments, err := p.queries.GetRootAndFirstDepthComments(ctx, threadId)

	if err != nil {
		return nil, err
	}

	commentEnts := []entities.Comment{}
	for _, comment := range comments {
		var deletedAt *time.Time
		if comment.DeletedAt.Valid {
			deletedAt = &comment.DeletedAt.Time
		}

		entitie := entities.
			NewComment().
			SetId(comment.ID).
			SetThreadId(comment.ThreadID).
			SetAddress(comment.Address).
			SetContent(comment.Content).
			SetVotes(comment.Votes).
			SetCreatedAt(comment.CreatedAt).
			SetDeletedAt(deletedAt).
			SetIsDeleted(comment.IsDeleted)
		commentEnts = append(commentEnts, entitie)
	}
	return commentEnts, nil
}

func (p *PostgresGateway) DeleteComment(ctx context.Context, id int32) error {
	_, err := p.queries.DeleteComment(ctx, id)

	if errors.Is(err, sql.ErrNoRows) {
		return common.ErrNotFound
	}

	return err
}

func (p *PostgresGateway) UpVoteComment(ctx context.Context, id int32, address string) error {
	return p.queries.CreateCommentUpVote(ctx, bindings.CreateCommentUpVoteParams{
		CommentID: id,
		Address:   address,
	})
}

func (p *PostgresGateway) DownVoteComment(ctx context.Context, id int32, address string) error {
	return p.queries.CreateCommentDownVote(ctx, bindings.CreateCommentDownVoteParams{
		CommentID: id,
		Address:   address,
	})
}

func (p *PostgresGateway) UnVoteComment(ctx context.Context, id int32, address string) error {
	return p.queries.CreateCommentUnVote(ctx, bindings.CreateCommentUnVoteParams{
		CommentID: id,
		Address:   address,
	})
}
