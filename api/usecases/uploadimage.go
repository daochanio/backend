package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
	"github.com/google/uuid"
)

type UploadImageUsecase struct {
	logger       common.Logger
	imageGateway ImageGateway
}

func NewUploadImageUsecase(logger common.Logger, imageGateway ImageGateway) *UploadImageUsecase {
	return &UploadImageUsecase{
		logger,
		imageGateway,
	}
}

type UploadImageInput struct {
	ContentType string  `validate:"oneof=image/jpeg image/png image/gif image/webp"`
	Bytes       *[]byte `validate:"required"`
}

func (u *UploadImageUsecase) Execute(ctx context.Context, input UploadImageInput) (entities.Image, error) {
	if err := common.ValidateStruct(input); err != nil {
		return entities.Image{}, err
	}

	id := uuid.New().String()
	ext := strings.Split(input.ContentType, "/")

	if len(ext) != 2 {
		return entities.Image{}, errors.New("invalid content type")
	}

	fileName := fmt.Sprintf("%s.%s", id, ext[1])

	return u.imageGateway.UploadImage(ctx, fileName, input.ContentType, input.Bytes)
}
