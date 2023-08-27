package usecases

import (
	"context"
	"io"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/core/entities"
	"github.com/daochanio/backend/core/gateways"
)

type UploadImage struct {
	logger common.Logger
	images gateways.Images
}

func NewUploadImageUsecase(logger common.Logger, images gateways.Images) *UploadImage {
	return &UploadImage{
		logger,
		images,
	}
}

type UploadImageInput struct {
	Reader io.Reader `validate:"required"`
}

func (u *UploadImage) Execute(ctx context.Context, input UploadImageInput) (*entities.Image, error) {
	if err := common.ValidateStruct(input); err != nil {
		return &entities.Image{}, err
	}

	return u.images.UploadImage(ctx, input.Reader)
}
