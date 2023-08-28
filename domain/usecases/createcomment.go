package usecases

import (
	"context"
	"fmt"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/domain/gateways"
)

type CreateComment struct {
	database  gateways.Database
	images    gateways.Images
	validator common.Validator
}

func NewCreateCommentUseCase(database gateways.Database, images gateways.Images, validator common.Validator) *CreateComment {
	return &CreateComment{
		database,
		images,
		validator,
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
	if err := u.validator.ValidateStruct(input); err != nil {
		return entities.Comment{}, err
	}

	image, err := u.images.GetImageByFileName(ctx, input.ImageFileName)

	if err != nil {
		return entities.Comment{}, err
	}

	if image == nil {
		return entities.Comment{}, fmt.Errorf("image not found %w", common.ErrNotFound)
	}

	return u.database.CreateComment(ctx, input.ThreadId, input.Address, input.RepliedToCommentId, input.Content, image)
}
