package usecases

import (
	"context"
	"errors"

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

func (u *DeleteCommentUseCase) Execute(ctx context.Context, input DeleteCommentInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	comment, err := u.dbGateway.GetCommentById(ctx, input.Id)

	if err != nil {
		return err
	}

	if comment.Address() != input.DeleterAddress {
		return errors.New("comment does not belong to the user")
	}

	return u.dbGateway.DeleteComment(ctx, input.Id)
}
