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

func (p *postgresGateway) CreateComment(
	ctx context.Context,
	threadId int64,
	address string,
	repliedToCommentId *int64,
	content string,
	image *entities.Image,
) (entities.Comment, error) {
	rep := pgtype.Int8{
		Valid: repliedToCommentId != nil,
	}

	if rep.Valid {
		rep.Int64 = *repliedToCommentId
	}

	id, err := p.queries.CreateComment(ctx, bindings.CreateCommentParams{
		ThreadID:                  threadId,
		Address:                   address,
		RepliedToCommentID:        rep,
		Content:                   content,
		ImageFileName:             image.FileName(),
		ImageOriginalUrl:          image.OriginalURL(),
		ImageOriginalContentType:  image.OriginalContentType(),
		ImageFormattedUrl:         image.FormattedURL(),
		ImageFormattedContentType: image.FormattedContentType(),
	})

	if err != nil {
		return entities.Comment{}, err
	}

	return p.GetCommentById(ctx, id)
}

func (p *postgresGateway) GetComments(ctx context.Context, threadId int64, offset int64, limit int64) ([]entities.Comment, int64, error) {
	dbComments, err := p.queries.GetComments(ctx, bindings.GetCommentsParams{
		ThreadID: threadId,
		Column2:  offset,
		Column3:  limit,
	})

	if err != nil {
		return nil, -1, err
	}

	count := int64(0)
	comments := []entities.Comment{}
	for _, dbComment := range dbComments {
		count = dbComment.FullCount

		var deletedAt *time.Time
		if dbComment.DeletedAt.Valid {
			deletedAt = &dbComment.DeletedAt.Time
		}
		image := entities.NewImage(dbComment.ImageFileName, dbComment.ImageOriginalUrl, dbComment.ImageOriginalContentType, dbComment.ImageFormattedUrl, dbComment.ImageFormattedContentType)
		user := toUser(
			dbComment.Address,
			dbComment.EnsName,
			dbComment.EnsAvatarFileName,
			dbComment.EnsAvatarOriginalUrl,
			dbComment.EnsAvatarOriginalContentType,
			dbComment.EnsAvatarFormattedUrl,
			dbComment.EnsAvatarFormattedContentType,
			dbComment.Reputation,
			dbComment.UserCreatedAt,
			dbComment.UserUpdatedAt,
		)
		repliedToComment := toRepliedToComment(
			dbComment.RID,
			dbComment.RContent,
			dbComment.RImageFileName,
			dbComment.RImageOriginalUrl,
			dbComment.RImageOriginalContentType,
			dbComment.RImageFormattedUrl,
			dbComment.RImageFormattedContentType,
			dbComment.RIsDeleted,
			dbComment.RCreatedAt,
			dbComment.RDeletedAt,
		)
		comment := entities.NewComment(entities.CommentParams{
			Id:               dbComment.ID,
			ThreadId:         dbComment.ThreadID,
			Content:          dbComment.Content,
			Image:            image,
			User:             user,
			RepliedToComment: repliedToComment,
			IsDeleted:        dbComment.IsDeleted,
			CreatedAt:        dbComment.CreatedAt.Time,
			DeletedAt:        deletedAt,
			Votes:            dbComment.Votes,
		})

		comments = append(comments, comment)
	}
	return comments, count, nil
}

func (p *postgresGateway) GetCommentById(ctx context.Context, id int64) (entities.Comment, error) {
	dbComment, err := p.queries.GetComment(ctx, id)

	if errors.Is(err, pgx.ErrNoRows) {
		return entities.Comment{}, common.ErrNotFound
	}

	if err != nil {
		return entities.Comment{}, err
	}

	var deletedAt *time.Time
	if dbComment.DeletedAt.Valid {
		deletedAt = &dbComment.DeletedAt.Time
	}

	image := entities.NewImage(
		dbComment.ImageFileName,
		dbComment.ImageOriginalUrl,
		dbComment.ImageOriginalContentType,
		dbComment.ImageFormattedUrl,
		dbComment.ImageFormattedContentType,
	)
	user := toUser(
		dbComment.Address,
		dbComment.EnsName,
		dbComment.EnsAvatarFileName,
		dbComment.EnsAvatarOriginalUrl,
		dbComment.EnsAvatarOriginalContentType,
		dbComment.EnsAvatarFormattedUrl,
		dbComment.EnsAvatarFormattedContentType,
		dbComment.Reputation,
		dbComment.UserCreatedAt,
		dbComment.UserUpdatedAt,
	)
	repliedToComment := toRepliedToComment(
		dbComment.RID,
		dbComment.RContent,
		dbComment.RImageFileName,
		dbComment.RImageOriginalUrl,
		dbComment.RImageOriginalContentType,
		dbComment.RImageFormattedUrl,
		dbComment.RImageFormattedContentType,
		dbComment.RIsDeleted,
		dbComment.RCreatedAt,
		dbComment.RDeletedAt,
	)
	entitie := entities.NewComment(entities.CommentParams{
		Id:               dbComment.ID,
		ThreadId:         dbComment.ThreadID,
		Content:          dbComment.Content,
		Image:            image,
		User:             user,
		RepliedToComment: repliedToComment,
		IsDeleted:        dbComment.IsDeleted,
		CreatedAt:        dbComment.CreatedAt.Time,
		DeletedAt:        deletedAt,
		Votes:            dbComment.Votes,
	})

	return entitie, nil
}

func (p *postgresGateway) DeleteComment(ctx context.Context, id int64) error {
	_, err := p.queries.DeleteComment(ctx, id)

	if errors.Is(err, pgx.ErrNoRows) {
		return common.ErrNotFound
	}

	return err
}

// assumes nil based on id validity
func toRepliedToComment(
	id pgtype.Int8,
	content pgtype.Text,
	imageFileName pgtype.Text,
	imageOriginalUrl pgtype.Text,
	imageOriginalContentType pgtype.Text,
	imageFormattedUrl pgtype.Text,
	imageFormattedContentType pgtype.Text,
	isDeleted pgtype.Bool,
	deletedAt pgtype.Timestamp,
	createdAt pgtype.Timestamp,
) *entities.Comment {
	if !id.Valid {
		return nil
	}
	var commentDeletedAt *time.Time
	if deletedAt.Valid {
		commentDeletedAt = &deletedAt.Time
	}
	repliedToImage := entities.NewImage(
		imageFileName.String,
		imageOriginalUrl.String,
		imageOriginalContentType.String,
		imageFormattedUrl.String,
		imageFormattedContentType.String,
	)
	comment := entities.NewComment(entities.CommentParams{
		Id:        id.Int64,
		Content:   content.String,
		Image:     repliedToImage,
		IsDeleted: isDeleted.Bool,
		CreatedAt: createdAt.Time,
		DeletedAt: commentDeletedAt,
	})
	return &comment
}
