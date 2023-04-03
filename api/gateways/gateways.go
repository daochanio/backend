package gateways

import (
	"context"

	"github.com/daochanio/backend/api/entities"
)

type IDatabaseGateway interface {
	CreateOrUpdateUser(ctx context.Context, address string, endName *string) error

	CreateThread(ctx context.Context, address string, content string) (int64, error)
	GetThreads(ctx context.Context, offset int32, limit int32) ([]entities.Thread, error)
	GetThreadById(ctx context.Context, id int64) (entities.Thread, error)
	DeleteThread(ctx context.Context, id int64) error
	UpVoteThread(ctx context.Context, id int64, address string) error
	DownVoteThread(ctx context.Context, id int64, address string) error
	UnVoteThread(ctx context.Context, id int64, address string) error

	CreateComment(ctx context.Context, threadId int64, address string, repliedToCommentId *int64, content string) (int64, error)
	GetComments(ctx context.Context, threadId int64, offset int32, limit int32) ([]entities.Comment, error)
	DeleteComment(ctx context.Context, id int64) error
	UpVoteComment(ctx context.Context, id int64, address string) error
	DownVoteComment(ctx context.Context, id int64, address string) error
	UnVoteComment(ctx context.Context, id int64, address string) error
}
