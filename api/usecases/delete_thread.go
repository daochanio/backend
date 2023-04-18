package usecases

import (
	"context"
	"errors"

	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

type DeleteThreadUseCase struct {
	dbGateway gateways.DatabaseGateway
}

func NewDeleteThreadUseCase(dbGateway gateways.DatabaseGateway) *DeleteThreadUseCase {
	return &DeleteThreadUseCase{
		dbGateway,
	}
}

type DeleteThreadInput struct {
	ThreadId       int64  `validate:"gt=0"`
	DeleterAddress string `validate:"eth_addr"`
}

func (u *DeleteThreadUseCase) Execute(ctx context.Context, input DeleteThreadInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	thread, err := u.dbGateway.GetThreadById(ctx, input.ThreadId)

	if err != nil {
		return err
	}

	if thread.Address() != input.DeleterAddress {
		return errors.New("thread does not belong to the user")
	}

	return u.dbGateway.DeleteThread(ctx, input.ThreadId)
}
