package gateways

import (
	"context"

	"github.com/daochanio/backend/api/entities"
)

type DatabaseGateway interface {
	CreateOrUpdateUser(ctx context.Context, address string, ensName *string) (entities.User, error)

	CreateThread(ctx context.Context, address string, content string) (int64, error)
	GetThreads(ctx context.Context, limit int32) ([]entities.Thread, error)
	GetThreadById(ctx context.Context, id int64) (entities.Thread, error)
	DeleteThread(ctx context.Context, id int64) error
	UpVoteThread(ctx context.Context, id int64, address string) error
	DownVoteThread(ctx context.Context, id int64, address string) error
	UnVoteThread(ctx context.Context, id int64, address string) error

	CreateComment(ctx context.Context, threadId int64, address string, repliedToCommentId *int64, content string) (int64, error)
	GetComments(ctx context.Context, threadId int64, offset int32, limit int32) ([]entities.Comment, error)
	GetCommentById(ctx context.Context, id int64) (entities.Comment, error)
	DeleteComment(ctx context.Context, id int64) error
	UpVoteComment(ctx context.Context, id int64, address string) error
	DownVoteComment(ctx context.Context, id int64, address string) error
	UnVoteComment(ctx context.Context, id int64, address string) error
}

type CacheGateway interface {
	GetChallengeByAddress(ctx context.Context, address string) (entities.Challenge, error)
	SaveChallenge(ctx context.Context, challenge entities.Challenge) error

	VerifyRateLimit(ctx context.Context, ipAddress string) error
}

type BlockchainGateway interface {
	GetENSNameFromAddress(ctx context.Context, address string) (string, error)
}
