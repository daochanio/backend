package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type CreateComment struct {
	database Database
	images   Images
}

func NewCreateCommentUseCase(database Database, images Images) *CreateComment {
	return &CreateComment{
		database,
		images,
	}
}

type CreateCommentInput struct {
	ThreadId           int64  `validate:"gt=0"`
	RepliedToCommentId *int64 `validate:"omitempty,gt=0"`
	Address            string `validate:"eth_addr"`
	Content            string `validate:"max=1000"`
	ImageFileName      string `validate:"max=100"`
}

func (u *CreateComment) Execute(ctx context.Context, input CreateCommentInput) (entities.Comment, error) {
	if err := common.ValidateStruct(input); err != nil {
		return entities.Comment{}, err
	}

	image, err := u.images.GetImageByFileName(ctx, input.ImageFileName)

	if err != nil {
		return entities.Comment{}, err
	}

	comment := entities.NewComment(entities.CommentParams{
		ThreadId: input.ThreadId,
		Address:  input.Address,
		Content:  input.Content,
		Image:    image,
	})

	return u.database.CreateComment(ctx, comment, input.RepliedToCommentId)
}
