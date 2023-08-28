package usecases

import (
	"context"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/domain/gateways"
)

type GetComments struct {
	validator common.Validator
	database  gateways.Database
}

func NewGetCommentsUseCase(validator common.Validator, database gateways.Database) *GetComments {
	return &GetComments{
		validator,
		database,
	}
}

type GetCommentsInput struct {
	ThreadId int64 `validate:"gt=0"`
	Offset   int64 `validate:"gte=0"`
	Limit    int64 `validate:"gt=0,lte=100"`
}

func (u *GetComments) Execute(ctx context.Context, input GetCommentsInput) ([]entities.Comment, int64, error) {
	if err := u.validator.ValidateStruct(input); err != nil {
		return nil, -1, err
	}

	return u.database.GetComments(ctx, input.ThreadId, input.Offset, input.Limit)
}
