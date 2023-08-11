package usecases

import (
	"context"
	"fmt"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type CreateThread struct {
	logger   common.Logger
	images   Images
	database Database
}

func NewCreateThreadUseCase(logger common.Logger, images Images, database Database) *CreateThread {
	return &CreateThread{
		logger,
		images,
		database,
	}
}

type CreateThreadInput struct {
	Address       string `validate:"eth_addr"`
	Title         string `validate:"max=100"`
	Content       string `validate:"max=1000"`
	ImageFileName string `validate:"max=100"`
}

func (u *CreateThread) Execute(ctx context.Context, input CreateThreadInput) (entities.Thread, error) {
	if err := common.ValidateStruct(input); err != nil {
		return entities.Thread{}, err
	}

	image, err := u.images.GetImageByFileName(ctx, input.ImageFileName)

	if err != nil {
		return entities.Thread{}, err
	}

	if image == nil {
		return entities.Thread{}, fmt.Errorf("image not found %w", common.ErrNotFound)
	}

	return u.database.CreateThread(ctx, input.Address, input.Title, input.Content, image)
}
