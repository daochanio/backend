package gateways

import (
	"context"
	"time"

	"github.com/daochanio/backend/api/entities"
)

type DatabaseGateway interface {
	CreateOrUpdateUser(ctx context.Context, address string, ensName *string) (entities.User, error)

	CreateThread(ctx context.Context, thread entities.Thread) (entities.Thread, error)
	GetThreads(ctx context.Context, limit int64) ([]entities.Thread, error)
	GetThreadById(ctx context.Context, id int64) (entities.Thread, error)
	DeleteThread(ctx context.Context, id int64) error
	UpVoteThread(ctx context.Context, id int64, address string) error
	DownVoteThread(ctx context.Context, id int64, address string) error
	UnVoteThread(ctx context.Context, id int64, address string) error

	CreateComment(ctx context.Context, comment entities.Comment, repliedToCommentId *int64) (entities.Comment, error)
	GetComments(ctx context.Context, threadId int64, offset int64, limit int64) ([]entities.Comment, int64, error)
	GetCommentById(ctx context.Context, id int64) (entities.Comment, error)
	DeleteComment(ctx context.Context, id int64) error
	UpVoteComment(ctx context.Context, id int64, address string) error
	DownVoteComment(ctx context.Context, id int64, address string) error
	UnVoteComment(ctx context.Context, id int64, address string) error
}

type CacheGateway interface {
	GetChallengeByAddress(ctx context.Context, address string) (entities.Challenge, error)
	SaveChallenge(ctx context.Context, challenge entities.Challenge) error

	VerifyRateLimit(ctx context.Context, key string, rate int, period time.Duration) error
}

type BlockchainGateway interface {
	GetENSNameFromAddress(ctx context.Context, address string) (string, error)
}

type ImageGateway interface {
	UploadImage(ctx context.Context, fileName string, contentType string, data *[]byte) (entities.Image, error)
	GetImageByFileName(ctx context.Context, fileName string) (entities.Image, error)
}
