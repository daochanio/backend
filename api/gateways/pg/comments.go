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

func (p *PostgresGateway) CreateComment(ctx context.Context, threadId int64, address string, repliedToCommentId *int64, content string) (int64, error) {
	rep := sql.NullInt64{
		Valid: repliedToCommentId != nil,
	}

	if rep.Valid {
		rep.Int64 = *repliedToCommentId
	}

	return p.queries.CreateComment(ctx, bindings.CreateCommentParams{
		ThreadID:           threadId,
		Address:            address,
		RepliedToCommentID: rep,
		Content:            content,
	})
}

func (p *PostgresGateway) GetComments(ctx context.Context, threadId int64, offset int32, limit int32) ([]entities.Comment, error) {
	comments, err := p.queries.GetComments(ctx, bindings.GetCommentsParams{
		ThreadID: threadId,
		Offset:   offset,
		Limit:    limit,
	})

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

		// set replying comment if exists
		if comment.RID.Valid {
			repliedToComment := entities.
				NewComment().
				SetId(comment.RID.Int64).
				SetAddress(comment.RAddress.String).
				SetContent(comment.RContent.String).
				SetCreatedAt(comment.RCreatedAt.Time).
				SetIsDeleted(comment.RIsDeleted.Bool)
			if comment.RDeletedAt.Valid {
				repliedToComment.SetDeletedAt(&comment.RDeletedAt.Time)
			}
			entitie = entitie.SetRepliedToComment(&repliedToComment)
		}

		commentEnts = append(commentEnts, entitie)
	}
	return commentEnts, nil
}

func (p *PostgresGateway) DeleteComment(ctx context.Context, id int64) error {
	_, err := p.queries.DeleteComment(ctx, id)

	if errors.Is(err, sql.ErrNoRows) {
		return common.ErrNotFound
	}

	return err
}

func (p *PostgresGateway) UpVoteComment(ctx context.Context, id int64, address string) error {
	return p.queries.CreateCommentUpVote(ctx, bindings.CreateCommentUpVoteParams{
		CommentID: id,
		Address:   address,
	})
}

func (p *PostgresGateway) DownVoteComment(ctx context.Context, id int64, address string) error {
	return p.queries.CreateCommentDownVote(ctx, bindings.CreateCommentDownVoteParams{
		CommentID: id,
		Address:   address,
	})
}

func (p *PostgresGateway) UnVoteComment(ctx context.Context, id int64, address string) error {
	return p.queries.CreateCommentUnVote(ctx, bindings.CreateCommentUnVoteParams{
		CommentID: id,
		Address:   address,
	})
}
