package usecases

import (
	"context"
	"fmt"
	"sync"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type GetThread struct {
	database Database
}

func NewGetThreadUseCase(database Database) *GetThread {
	return &GetThread{
		database,
	}
}

type GetThreadInput struct {
	ThreadId      int64 `validate:"gt=0"`
	CommentOffset int64 `validate:"gte=0"`
	CommentLimit  int64 `validate:"gt=0,lte=100"`
}

func (u *GetThread) Execute(ctx context.Context, input GetThreadInput) (entities.Thread, int64, error) {
	if err := common.ValidateStruct(input); err != nil {
		return entities.Thread{}, -1, fmt.Errorf("invalid input: %w", err)
	}

	// threads and comments can be fetched concurrently
	var thread entities.Thread
	var comments []entities.Comment
	var commentsCount int64
	var threadErr error
	var commentsErr error

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		thread, threadErr = u.database.GetThreadById(ctx, input.ThreadId)
	}()
	go func() {
		defer wg.Done()
		comments, commentsCount, commentsErr = u.database.GetComments(ctx, input.ThreadId, input.CommentOffset, input.CommentLimit)
	}()
	wg.Wait()

	if threadErr != nil {
		return entities.Thread{}, -1, fmt.Errorf("failed to fetch thread: %w", threadErr)
	}

	if commentsErr != nil {
		return entities.Thread{}, -1, fmt.Errorf("failed to fetch comments: %w", commentsErr)
	}

	thread.SetComments(&comments)

	return thread, commentsCount, nil
}
