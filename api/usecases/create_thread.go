package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

type CreateThreadUseCase struct {
	logger       common.Logger
	imageGateway gateways.ImageGateway
	dbGateway    gateways.DatabaseGateway
}

func NewCreateThreadUseCase(logger common.Logger, imageGateway gateways.ImageGateway, dbGateway gateways.DatabaseGateway) *CreateThreadUseCase {
	return &CreateThreadUseCase{
		logger,
		imageGateway,
		dbGateway,
	}
}

type CreateThreadInput struct {
	Address       string `validate:"eth_addr"`
	Title         string `validate:"max=100"`
	Content       string `validate:"max=1000"`
	ImageFileName string `validate:"max=100"`
}

func (u *CreateThreadUseCase) Execute(ctx context.Context, input CreateThreadInput) (entities.Thread, error) {
	if err := common.ValidateStruct(input); err != nil {
		return entities.Thread{}, err
	}

	image, err := u.imageGateway.GetImageByFileName(ctx, input.ImageFileName)

	if err != nil {
		return entities.Thread{}, err
	}

	thread := entities.NewThread(entities.ThreadParams{
		Address: input.Address,
		Title:   input.Title,
		Content: input.Content,
		Image:   image,
	})

	return u.dbGateway.CreateThread(ctx, thread)
}
