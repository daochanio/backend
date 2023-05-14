package usecases

import (
	"context"
	"errors"

	"github.com/daochanio/backend/common"
)

type DeleteThread struct {
	database Database
}

func NewDeleteThreadUseCase(database Database) *DeleteThread {
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

	if thread.Address() != input.DeleterAddress {
		return errors.New("thread does not belong to the user")
	}

	return u.database.DeleteThread(ctx, input.ThreadId)
}
