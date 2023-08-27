package usecases

import (
	"context"
	"errors"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/core/gateways"
)

type DeleteThread struct {
	database gateways.Database
}

func NewDeleteThreadUseCase(database gateways.Database) *DeleteThread {
	return &DeleteThread{
		database,
	}
}

type DeleteThreadInput struct {
	ThreadId       int64  `validate:"gt=0"`
	DeleterAddress string `validate:"eth_addr"`
}

func (u *DeleteThread) Execute(ctx context.Context, input DeleteThreadInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	thread, err := u.database.GetThreadById(ctx, input.ThreadId)

	if err != nil {
		return err
	}

	user := thread.User()
	if user.Address() != input.DeleterAddress {
		return errors.New("thread does not belong to the user")
	}

	return u.database.DeleteThread(ctx, input.ThreadId)
}
