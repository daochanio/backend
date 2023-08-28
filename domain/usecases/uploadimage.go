package usecases

import (
	"context"
	"io"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/domain/gateways"
)

type UploadImage struct {
	logger    common.Logger
	validator common.Validator
	images    gateways.Images
}

func NewUploadImageUsecase(logger common.Logger, validator common.Validator, images gateways.Images) *UploadImage {
	return &UploadImage{
		logger,
		validator,
		images,
	}
}

type UploadImageInput struct {
	Reader io.Reader `validate:"required"`
}

func (u *UploadImage) Execute(ctx context.Context, input UploadImageInput) (*entities.Image, error) {
	if err := u.validator.ValidateStruct(input); err != nil {
		return &entities.Image{}, err
	}

	return u.images.UploadImage(ctx, input.Reader)
}
