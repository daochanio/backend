package usecases

import (
	"context"
	"fmt"

	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

type CreateCommentVoteUseCase struct {
	dbGateway gateways.DatabaseGateway
}

func NewCreateCommentVoteUseCase(dbGateway gateways.DatabaseGateway) *CreateCommentVoteUseCase {
	return &CreateCommentVoteUseCase{
		dbGateway,
	}
}

type CreateCommentVoteInput struct {
	CommentId int64           `validate:"gt=0"`
	Address   string          `validate:"eth_addr"`
	Vote      common.VoteType `validate:"oneof=up down undo"`
}

func (u *CreateCommentVoteUseCase) Execute(ctx context.Context, input CreateCommentVoteInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	switch input.Vote {
	case common.UpVote:
		return u.dbGateway.UpVoteComment(ctx, input.CommentId, input.Address)
	case common.DownVote:
		return u.dbGateway.DownVoteComment(ctx, input.CommentId, input.Address)
	case common.UnVote:
		return u.dbGateway.UnVoteComment(ctx, input.CommentId, input.Address)
	default:
		return fmt.Errorf("invalid vote type: %v", input.Vote)
	}
}
