package usecases

import (
	"context"
	"fmt"

	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

type VoteThreadUseCase struct {
	dbGateway gateways.IDatabaseGateway
}

func NewVoteThreadUseCase(dbGateway gateways.IDatabaseGateway) *VoteThreadUseCase {
	return &VoteThreadUseCase{
		dbGateway,
	}
}

type VoteType string

const (
	UpVote   VoteType = "up"
	DownVote VoteType = "down"
	UnVote   VoteType = "undo"
)

type VoteThreadInput struct {
	ThreadId int32
	Address  string   `validate:"eth_addr"`
	Vote     VoteType `validate:"oneof=up down undo"`
}

func (u *VoteThreadUseCase) Execute(ctx context.Context, input VoteThreadInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	switch input.Vote {
	case UpVote:
		return u.dbGateway.UpVoteThread(ctx, input.ThreadId, input.Address)
	case DownVote:
		return u.dbGateway.DownVoteThread(ctx, input.ThreadId, input.Address)
	case UnVote:
		return u.dbGateway.UnVoteThread(ctx, input.ThreadId, input.Address)
	default:
		return fmt.Errorf("invalid vote type: %v", input.Vote)
	}
}
