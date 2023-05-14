package usecases

import (
	"context"
	"errors"

	"github.com/daochanio/backend/common"
)

type DeleteComment struct {
	database Database
}

func NewDeleteCommentUseCase(database Database) *DeleteComment {
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

	if comment.Address() != input.DeleterAddress {
		return errors.New("comment does not belong to the user")
	}

	return u.database.DeleteComment(ctx, input.Id)
}
