package usecases

import (
	"context"

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

func (u *CreateThreadUseCase) Execute(ctx context.Context, input CreateThreadInput) (int64, error) {
	if err := common.ValidateStruct(input); err != nil {
		return 0, err
	}

	image, err := u.imageGateway.GetImageById(ctx, input.ImageFileName)

	if err != nil {
		return 0, err
	}

	return u.dbGateway.CreateThread(ctx, input.Address, input.Title, input.Content, image.FileName(), image.Url(), image.ContentType())
}
