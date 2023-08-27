package usecases

import (
	"context"
	"errors"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/gateways"
)

type DeleteComment struct {
	database gateways.Database
}

func NewDeleteCommentUseCase(database gateways.Database) *DeleteComment {
	return &DeleteComment{
		database,
	}
}

type DeleteCommentInput struct {
	Id             int64  `validate:"gt=0"`
	DeleterAddress string `validate:"eth_addr"`
}

func (u *DeleteComment) Execute(ctx context.Context, input DeleteCommentInput) error {
	if err := common.ValidateStruct(input); err != nil {
		return err
	}

	comment, err := u.database.GetCommentById(ctx, input.Id)

	if err != nil {
		return err
	}

	user := comment.User()
	if user.Address() != input.DeleterAddress {
		return errors.New("comment does not belong to the user")
	}

	return u.database.DeleteComment(ctx, input.Id)
}
