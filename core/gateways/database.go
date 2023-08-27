package gateways

import (
	"context"
	"math/big"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/core/entities"
)

type DatabaseConfig struct {
	ConnectionString string
	MinConnections   int32
	MaxConnections   int32
}

type Database interface {
	Start(ctx context.Context, config DatabaseConfig)
	Shutdown(ctx context.Context)

	GetChallengeByAddress(ctx context.Context, address string) (entities.Challenge, error)
	SaveChallenge(ctx context.Context, challenge entities.Challenge) error

	GetUserByAddress(ctx context.Context, address string) (entities.User, error)
	GetThreads(ctx context.Context, limit int64) ([]entities.Thread, error)
	GetThreadById(ctx context.Context, threadId int64) (entities.Thread, error)
	GetComments(ctx context.Context, threadId int64, offset int64, limit int64) ([]entities.Comment, int64, error)
	GetCommentById(ctx context.Context, commentId int64) (entities.Comment, error)

	UpsertUser(ctx context.Context, address string) error
	UpdateUser(ctx context.Context, address string, name *string, avatar *entities.Image) error
	CreateComment(ctx context.Context, threadId int64, address string, repliedToCommentId *int64, content string, image *entities.Image) (entities.Comment, error)
	CreateThread(ctx context.Context, address string, title string, content string, image *entities.Image) (entities.Thread, error)
	CreateVote(ctx context.Context, vote entities.Vote) error
	DeleteThread(ctx context.Context, threadId int64) error
	DeleteComment(ctx context.Context, commentId int64) error
	AggregateVotes(ctx context.Context, id int64, voteType common.VoteType) error

	GetLastIndexedBlock(ctx context.Context) (*big.Int, error)
	UpdateLastIndexedBlock(ctx context.Context, block *big.Int) error
	InsertTransferEvents(ctx context.Context, from *big.Int, to *big.Int, transfers []entities.Transfer) error
	UpdateReputation(ctx context.Context, addresses []string) error
}
