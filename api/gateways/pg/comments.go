package pg

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/db/bindings"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p *postgresGateway) CreateComment(ctx context.Context, comment entities.Comment, repliedToCommentId *int64) (entities.Comment, error) {
	rep := pgtype.Int8{
		Valid: repliedToCommentId != nil,
	}

	if rep.Valid {
		rep.Int64 = *repliedToCommentId
	}

	createdComment, err := p.queries.CreateComment(ctx, bindings.CreateCommentParams{
		ThreadID:           comment.ThreadId(),
		Address:            comment.Address(),
		RepliedToCommentID: rep,
		Content:            comment.Content(),
		ImageFileName:      comment.Image().FileName(),
		ImageUrl:           comment.Image().Url(),
		ImageContentType:   comment.Image().ContentType(),
	})

	if err != nil {
		return entities.Comment{}, err
	}

	var deletedAt *time.Time
	if createdComment.DeletedAt.Valid {
		deletedAt = &createdComment.DeletedAt.Time
	}

	image := entities.NewImage(createdComment.ImageFileName, createdComment.ImageUrl, createdComment.ImageContentType)
	ent := entities.NewComment(entities.CommentParams{
		Id:        createdComment.ID,
		ThreadId:  createdComment.ThreadID,
		Address:   createdComment.Address,
		Content:   createdComment.Content,
		Image:     image,
		IsDeleted: createdComment.IsDeleted,
		CreatedAt: createdComment.CreatedAt.Time,
		DeletedAt: deletedAt,
	})

	// set replying comment if exists
	if createdComment.RepliedToCommentID.Valid {
		repliedToComment, err := p.GetCommentById(ctx, createdComment.RepliedToCommentID.Int64)

		if err != nil {
			p.logger.Warn(ctx).Err(err).Msg("failed to get replied to comment")
			return ent, nil
		}

		ent.SetRepliedToComment(&repliedToComment)
	}

	return ent, nil
}

func (p *postgresGateway) GetComments(ctx context.Context, threadId int64, offset int64, limit int64) ([]entities.Comment, int64, error) {
	comments, err := p.queries.GetComments(ctx, bindings.GetCommentsParams{
		ThreadID: threadId,
		Column2:  offset,
		Column3:  limit,
	})

	if err != nil {
		return nil, -1, err
	}

	count := int64(0)
	commentEnts := []entities.Comment{}
	for _, comment := range comments {
		count = comment.FullCount

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
			CreatedAt:        comment.CreatedAt.Time,
			DeletedAt:        deletedAt,
			Votes:            comment.Votes,
		})

		commentEnts = append(commentEnts, entitie)
	}
	return commentEnts, count, nil
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
		CreatedAt: comment.CreatedAt.Time,
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

func (p *postgresGateway) CreateCommentVote(ctx context.Context, vote entities.Vote) error {
	switch vote.Value() {
	case common.Upvote:
		return p.queries.CreateCommentUpVote(ctx, bindings.CreateCommentUpVoteParams{
			CommentID: vote.Id(),
			Address:   vote.Address(),
		})
	case common.Downvote:
		return p.queries.CreateCommentDownVote(ctx, bindings.CreateCommentDownVoteParams{
			CommentID: vote.Id(),
			Address:   vote.Address(),
		})
	case common.Unvote:
		return p.queries.CreateCommentUnVote(ctx, bindings.CreateCommentUnVoteParams{
			CommentID: vote.Id(),
			Address:   vote.Address(),
		})
	default:
		return errors.New("invalid vote value")
	}
}

func (p *postgresGateway) AggregateCommentVotes(ctx context.Context, commentId int64) error {
	return p.queries.AggregateCommentVotes(ctx, commentId)
}
