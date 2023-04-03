package usecases

import (
	"context"

	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

type DeleteCommentUseCase struct {
	dbGateway gateways.IDatabaseGateway
}

func NewDeleteCommentUseCase(dbGateway gateways.IDatabaseGateway) *DeleteCommentUseCase {
	return &DeleteCommentUseCase{
		dbGateway,
	}
}

type DeleteCommentInput struct {
	Id             int64  `validate:"gt=0"`
	DeleterAddress string `validate:"eth_addr"`
}

// TODO: check if the comment belongs to the user (or is a moderate once supported) before deleting
func (u *DeleteCommentUseCase) Execute(ctx context.Context, input DeleteCommentInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	return u.dbGateway.DeleteComment(ctx, input.Id)
}
