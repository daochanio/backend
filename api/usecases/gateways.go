package usecases

import (
	"context"
	"time"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type Database interface {
	Start(ctx context.Context)
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
	CreateComment(ctx context.Context, threadId int64, address string, repliedToCommentId *int64, content string, imageFileName string, imageUrl string, imageContentType string) (entities.Comment, error)
	CreateThread(ctx context.Context, address string, title string, content string, imageFileName string, imageUrl string, imageContentType string) (entities.Thread, error)
	CreateVote(ctx context.Context, vote entities.Vote) error
	DeleteThread(ctx context.Context, threadId int64) error
	DeleteComment(ctx context.Context, commentId int64) error
	AggregateVotes(ctx context.Context, id int64, voteType common.VoteType) error
}

type Cache interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
	VerifyRateLimit(ctx context.Context, key string, rate int, period time.Duration) error
}

type Stream interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
	PublishSignin(ctx context.Context, address string) error
	PublishVote(ctx context.Context, vote entities.Vote) error
}

type Blockchain interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
	GetNameByAddress(ctx context.Context, address string) (*string, error)
	GetAvatarURIByName(ctx context.Context, name string) (*string, error)
	GetNFTURI(ctx context.Context, standard string, address string, id string) (string, error)
}

type Storage interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
	UploadImage(ctx context.Context, fileName string, contentType string, data *[]byte) (entities.Image, error)
	GetImageByFileName(ctx context.Context, fileName string) (*entities.Image, error)
}

type SafeProxy interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
	DownloadImage(ctx context.Context, uri string) (*[]byte, string, error)
	GetNFTImageURI(ctx context.Context, uri string) (string, error)
}
