package usecases

import (
	"context"

	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

type CreateCommentUseCase struct {
	dbGateway gateways.IDatabaseGateway
}

func NewCreateCommentUseCase(dbGateway gateways.IDatabaseGateway) *CreateCommentUseCase {
	return &CreateCommentUseCase{
		dbGateway,
	}
}

type CreateCommentInput struct {
	ThreadId        int64  `validate:"gt=0"`
	ParentCommentId *int64 `validate:"omitempty,gt=0"`
	Address         string `validate:"eth_addr"`
	Content         string `validate:"max=1000"`
}

func (u *CreateCommentUseCase) Execute(ctx context.Context, input CreateCommentInput) (int64, error) {
	if err := common.ValidateStruct(input); err != nil {
		return 0, err
	}

	return u.dbGateway.CreateComment(ctx, input.ThreadId, input.Address, input.ParentCommentId, input.Content)
}
