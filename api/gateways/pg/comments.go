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

func (p *postgresGateway) CreateComment(ctx context.Context, comment entities.Comment, repliedToCommentId *int64) (int64, error) {
	rep := sql.NullInt64{
		Valid: repliedToCommentId != nil,
	}

	if rep.Valid {
		rep.Int64 = *repliedToCommentId
	}

	return p.queries.CreateComment(ctx, bindings.CreateCommentParams{
		ThreadID:           comment.ThreadId(),
		Address:            comment.Address(),
		RepliedToCommentID: rep,
		Content:            comment.Content(),
		ImageFileName:      comment.Image().FileName(),
		ImageUrl:           comment.Image().Url(),
		ImageContentType:   comment.Image().ContentType(),
	})
}

func (p *postgresGateway) GetComments(ctx context.Context, threadId int64, offset int32, limit int32) ([]entities.Comment, error) {
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

		image := entities.NewImage(comment.ImageFileName, comment.ImageUrl, comment.ImageContentType)
		// set replying comment if exists
		var repliedToComment *entities.Comment
		if comment.RID.Valid {
			var deletedAt *time.Time
			if comment.RDeletedAt.Valid {
				deletedAt = &comment.RDeletedAt.Time
			}
			repliedToImage := entities.NewImage(comment.RImageFileName.String, comment.RImageUrl.String, comment.RImageContentType.String)
			comment := entities.NewComment(entities.CommentParams{
				Id:        comment.RID.Int64,
				Address:   comment.RAddress.String,
				Content:   comment.RContent.String,
				Image:     repliedToImage,
				IsDeleted: comment.RIsDeleted.Bool,
				CreatedAt: comment.RCreatedAt.Time,
				DeletedAt: deletedAt,
			})
			repliedToComment = &comment
		}

		entitie := entities.NewComment(entities.CommentParams{
			Id:               comment.ID,
			ThreadId:         comment.ThreadID,
			Address:          comment.Address,
			Content:          comment.Content,
			Image:            image,
			RepliedToComment: repliedToComment,
			IsDeleted:        comment.IsDeleted,
			CreatedAt:        comment.CreatedAt,
			DeletedAt:        deletedAt,
			Votes:            comment.Votes,
		})

		commentEnts = append(commentEnts, entitie)
	}
	return commentEnts, nil
}

// does not return with the hydrated replied to comment
func (p *postgresGateway) GetCommentById(ctx context.Context, id int64) (entities.Comment, error) {
	comment, err := p.queries.GetComment(ctx, id)

	if errors.Is(err, sql.ErrNoRows) {
		return entities.Comment{}, common.ErrNotFound
	}

	if err != nil {
		return entities.Comment{}, err
	}

	var deletedAt *time.Time
	if comment.DeletedAt.Valid {
		deletedAt = &comment.DeletedAt.Time
	}

	image := entities.NewImage(comment.ImageFileName, comment.ImageUrl, comment.ImageContentType)
	entitie := entities.NewComment(entities.CommentParams{
		Id:        comment.ID,
		ThreadId:  comment.ThreadID,
		Address:   comment.Address,
		Content:   comment.Content,
		Image:     image,
		IsDeleted: comment.IsDeleted,
		CreatedAt: comment.CreatedAt,
		DeletedAt: deletedAt,
		Votes:     comment.Votes,
	})

	return entitie, nil
}

func (p *postgresGateway) DeleteComment(ctx context.Context, id int64) error {
	_, err := p.queries.DeleteComment(ctx, id)

	if errors.Is(err, sql.ErrNoRows) {
		return common.ErrNotFound
	}

	return err
}

func (p *postgresGateway) UpVoteComment(ctx context.Context, id int64, address string) error {
	return p.queries.CreateCommentUpVote(ctx, bindings.CreateCommentUpVoteParams{
		CommentID: id,
		Address:   address,
	})
}

func (p *postgresGateway) DownVoteComment(ctx context.Context, id int64, address string) error {
	return p.queries.CreateCommentDownVote(ctx, bindings.CreateCommentDownVoteParams{
		CommentID: id,
		Address:   address,
	})
}

func (p *postgresGateway) UnVoteComment(ctx context.Context, id int64, address string) error {
	return p.queries.CreateCommentUnVote(ctx, bindings.CreateCommentUnVoteParams{
		CommentID: id,
		Address:   address,
	})
}
