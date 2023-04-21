package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

type CreateCommentUseCase struct {
	dbGateway    gateways.DatabaseGateway
	imageGateway gateways.ImageGateway
}

func NewCreateCommentUseCase(dbGateway gateways.DatabaseGateway, imageGateway gateways.ImageGateway) *CreateCommentUseCase {
	return &CreateCommentUseCase{
		dbGateway,
		imageGateway,
	}
}

type CreateCommentInput struct {
	ThreadId           int64  `validate:"gt=0"`
	RepliedToCommentId *int64 `validate:"omitempty,gt=0"`
	Address            string `validate:"eth_addr"`
	Content            string `validate:"max=1000"`
	ImageFileName      string `validate:"max=100"`
}

func (u *CreateCommentUseCase) Execute(ctx context.Context, input CreateCommentInput) (int64, error) {
	if err := common.ValidateStruct(input); err != nil {
		return 0, err
	}

	image, err := u.imageGateway.GetImageByFileName(ctx, input.ImageFileName)

	if err != nil {
		return 0, err
	}

	comment := entities.NewComment(entities.CommentParams{
		ThreadId: input.ThreadId,
		Address:  input.Address,
		Content:  input.Content,
		Image:    image,
	})

	return u.dbGateway.CreateComment(ctx, comment, input.RepliedToCommentId)
}
