package index

import (
	"context"
	"errors"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/core/usecases"
)

type Indexer interface {
	Start(ctx context.Context, config IndexerConfig)
	Shutdown(ctx context.Context)
}

type indexer struct {
	logger      common.Logger
	indexBlocks *usecases.IndexBlocks
}

func NewIndexer(logger common.Logger, indexBlocks *usecases.IndexBlocks) Indexer {
	return &indexer{
		logger,
		indexBlocks,
	}
}

type IndexerConfig struct {
	Interval      time.Duration
	MaxBlockRange int64
	ReorgOffset   int64
}

func (i *indexer) Start(ctx context.Context, config IndexerConfig) {
	i.logger.Info(ctx).Msg("starting indexer")

	for {
		select {
		case <-ctx.Done():
			i.logger.Info(ctx).Msg("indexer stopped")
			return
		default:
			if err := i.indexBlocks.Execute(ctx, usecases.IndexBlocksInput{
				MaxBlockRange: config.MaxBlockRange,
				ReorgOffset:   config.ReorgOffset,
			}); err != nil {
				if !errors.Is(err, common.ErrNoNewBlocks) {
					// error log if the error was anything other than no new blocks
					i.logger.Error(ctx).Err(err).Msg("could not index blocks")
				}

				// sleep a bit if we errored while indexing to avoid spamming the blockchain provider and allow time recover from transient errors or to simply wait for a new block
				time.Sleep(config.Interval)
			}
		}
	}
}

func (i *indexer) Shutdown(ctx context.Context) {
	i.logger.Info(ctx).Msg("shutting down indexer")
}
