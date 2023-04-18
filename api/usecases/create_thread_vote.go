package usecases

import (
	"context"
	"fmt"

	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

type CreateThreadVoteUseCase struct {
	dbGateway gateways.DatabaseGateway
}

func NewCreateThreadVoteUseCase(dbGateway gateways.DatabaseGateway) *CreateThreadVoteUseCase {
	return &CreateThreadVoteUseCase{
		dbGateway,
	}
}

type CreateThreadVoteInput struct {
	ThreadId int64           `validate:"gt=0"`
	Address  string          `validate:"eth_addr"`
	Vote     common.VoteType `validate:"oneof=up down undo"`
}

func (u *CreateThreadVoteUseCase) Execute(ctx context.Context, input CreateThreadVoteInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	switch input.Vote {
	case common.UpVote:
		return u.dbGateway.UpVoteThread(ctx, input.ThreadId, input.Address)
	case common.DownVote:
		return u.dbGateway.DownVoteThread(ctx, input.ThreadId, input.Address)
	case common.UnVote:
		return u.dbGateway.UnVoteThread(ctx, input.ThreadId, input.Address)
	default:
		return fmt.Errorf("invalid vote type: %v", input.Vote)
	}
}
