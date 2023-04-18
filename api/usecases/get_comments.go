package usecases

import (
	"context"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

type GetCommentsUseCase struct {
	dbGateway gateways.DatabaseGateway
}

func NewGetCommentsUseCase(dbGateway gateways.DatabaseGateway) *GetCommentsUseCase {
	return &GetCommentsUseCase{
		dbGateway,
	}
}

type GetCommentsInput struct {
	ThreadId int64 `validate:"gt=0"`
	Offset   int32 `validate:"gte=0"`
	Limit    int32 `validate:"gte=0"`
}

func (u *GetCommentsUseCase) Execute(ctx context.Context, input GetCommentsInput) ([]entities.Comment, error) {
	if err := common.ValidateStruct(input); err != nil {
		return nil, err
	}

	return u.dbGateway.GetComments(ctx, input.ThreadId, input.Offset, input.Limit)
}
