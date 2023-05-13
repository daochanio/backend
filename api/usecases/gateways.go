package usecases

import (
	"context"
	"time"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type DatabaseGateway interface {
	GetThreads(ctx context.Context, limit int64) ([]entities.Thread, error)
	GetThreadById(ctx context.Context, threadId int64) (entities.Thread, error)
	GetComments(ctx context.Context, threadId int64, offset int64, limit int64) ([]entities.Comment, int64, error)
	GetCommentById(ctx context.Context, commentId int64) (entities.Comment, error)

	UpsertUser(ctx context.Context, address string) error
	UpdateUser(ctx context.Context, address string, ensName *string) error
	CreateComment(ctx context.Context, comment entities.Comment, repliedToCommentId *int64) (entities.Comment, error)
	CreateThread(ctx context.Context, thread entities.Thread) (entities.Thread, error)
	CreateVote(ctx context.Context, vote entities.Vote) error
	DeleteThread(ctx context.Context, threadId int64) error
	DeleteComment(ctx context.Context, commentId int64) error
	AggregateVotes(ctx context.Context, id int64, voteType common.VoteType) error
}

type CacheGateway interface {
	GetChallengeByAddress(ctx context.Context, address string) (entities.Challenge, error)
	SaveChallenge(ctx context.Context, challenge entities.Challenge) error
	VerifyRateLimit(ctx context.Context, key string, rate int, period time.Duration) error
}

type MessageGateway interface {
	PublishSignin(ctx context.Context, address string) error
	PublishVote(ctx context.Context, vote entities.Vote) error
}

type BlockchainGateway interface {
	GetNameFromAddress(ctx context.Context, address string) (string, error)
}

type ImageGateway interface {
	UploadImage(ctx context.Context, fileName string, contentType string, data *[]byte) (entities.Image, error)
	GetImageByFileName(ctx context.Context, fileName string) (entities.Image, error)
}
