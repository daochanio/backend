package usecases

import (
	"context"

	"github.com/daochanio/backend/api/gateways"
)

type DeleteThreadUseCase struct {
	dbGateway gateways.IDatabaseGateway
}

func NewDeleteThreadUseCase(dbGateway gateways.IDatabaseGateway) *DeleteThreadUseCase {
	return &DeleteThreadUseCase{
		dbGateway,
	}
}

type DeleteThreadInput struct {
	ThreadId       int64  `validate:"gt=0"`
	DeleterAddress string `validate:"eth_addr"`
}

// TODO: only allow the thread creator (or in the future moderators) to delete the thread
func (u *DeleteThreadUseCase) Execute(ctx context.Context, input DeleteThreadInput) error {
	return u.dbGateway.DeleteThread(ctx, input.ThreadId)
}
