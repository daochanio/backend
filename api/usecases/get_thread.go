package usecases

import (
	"context"
	"fmt"
	"sync"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/gateways"
	"github.com/daochanio/backend/common"
)

type GetThreadUseCase struct {
	dbGateway gateways.DatabaseGateway
}

func NewGetThreadUseCase(dbGateway gateways.DatabaseGateway) *GetThreadUseCase {
	return &GetThreadUseCase{
		dbGateway,
	}
}

// get thread input
type GetThreadInput struct {
	ThreadId      int64 `validate:"gt=0"`
	CommentOffset int64 `validate:"gte=0"`
	CommentLimit  int64 `validate:"gt=0,lte=100"`
}

// Fetches a thread and its comments
func (u *GetThreadUseCase) Execute(ctx context.Context, input GetThreadInput) (entities.Thread, int64, error) {
	if err := common.ValidateStruct(input); err != nil {
		return entities.Thread{}, -1, fmt.Errorf("invalid input: %w", err)
	}

	// threads and comments can be fetched concurrently
	threadChan := make(chan entities.Thread, 1)
	commentsChan := make(chan []entities.Comment, 1)
	commentsCountChan := make(chan int64, 1)
	errs := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		thread, err := u.dbGateway.GetThreadById(ctx, input.ThreadId)
		if err != nil {
			errs <- fmt.Errorf("failed to get thread: %w", err)
			return
		}
		threadChan <- thread
	}()
	go func() {
		defer wg.Done()
		comments, count, err := u.dbGateway.GetComments(ctx, input.ThreadId, input.CommentOffset, input.CommentLimit)
		if err != nil {
			errs <- fmt.Errorf("failed to get comments: %w", err)
			return
		}
		commentsChan <- comments
		commentsCountChan <- count
	}()
	wg.Wait()

	close(threadChan)
	close(commentsChan)
	close(commentsCountChan)
	close(errs)

	for err := range errs {
		if err != nil {
			return entities.Thread{}, -1, err
		}
	}

	thread := <-threadChan
	comments := <-commentsChan
	count := <-commentsCountChan
	thread.SetComments(&comments)

	return thread, count, nil
}
