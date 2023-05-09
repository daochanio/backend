package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type GetCommentsUseCase struct {
	dbGateway DatabaseGateway
}

func NewGetCommentsUseCase(dbGateway DatabaseGateway) *GetCommentsUseCase {
	return &GetCommentsUseCase{
		dbGateway,
	}
}

type GetCommentsInput struct {
	ThreadId int64 `validate:"gt=0"`
	Offset   int64 `validate:"gte=0"`
	Limit    int64 `validate:"gt=0,lte=100"`
}

func (u *GetCommentsUseCase) Execute(ctx context.Context, input GetCommentsInput) ([]entities.Comment, int64, error) {
	if err := common.ValidateStruct(input); err != nil {
		return nil, -1, err
	}

	return u.dbGateway.GetComments(ctx, input.ThreadId, input.Offset, input.Limit)
}
