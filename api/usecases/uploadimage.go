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

type UploadImage struct {
	logger  common.Logger
	storage Storage
}

func NewUploadImageUsecase(logger common.Logger, storage Storage) *UploadImage {
	return &UploadImage{
		logger,
		storage,
	}
}

type UploadImageInput struct {
	ContentType string  `validate:"oneof=image/jpeg image/png image/gif image/webp"`
	Bytes       *[]byte `validate:"required"`
}

func (u *UploadImage) Execute(ctx context.Context, input UploadImageInput) (entities.Image, error) {
	if err := common.ValidateStruct(input); err != nil {
		return entities.Image{}, err
	}

	id := uuid.New().String()
	ext := strings.Split(input.ContentType, "/")

	if len(ext) != 2 {
		return entities.Image{}, errors.New("invalid content type")
	}

	fileName := fmt.Sprintf("%s.%s", id, ext[1])

	return u.storage.UploadImage(ctx, fileName, input.ContentType, input.Bytes)
}
