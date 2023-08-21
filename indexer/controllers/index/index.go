package index

import (
	"context"
	"errors"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/settings"
	"github.com/daochanio/backend/indexer/usecases"
)

type Indexer interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
}

type indexer struct {
	logger      common.Logger
	settings    settings.Settings
	indexBlocks *usecases.IndexBlocks
}

func NewIndexer(logger common.Logger, settings settings.Settings, indexBlocks *usecases.IndexBlocks) Indexer {
	return &indexer{
		logger,
		settings,
		indexBlocks,
	}
}

func (i *indexer) Start(ctx context.Context) {
	i.logger.Info(ctx).Msg("starting indexer")

	for {
		select {
		case <-ctx.Done():
			i.logger.Info(ctx).Msg("indexer stopped")
			return
		default:
			if err := i.indexBlocks.Execute(ctx); err != nil {
				if !errors.Is(err, common.ErrNoNewBlocks) {
					// error log if the error was anything other than no new blocks
					i.logger.Error(ctx).Err(err).Msg("could not index blocks")
				}

				// sleep a bit if we errored while indexing to avoid spamming the blockchain provider and allow time recover from transient errors or to simply wait for a new block
				time.Sleep(time.Duration(i.settings.IntervalSeconds()) * time.Second)
			}
		}
	}
}

func (i *indexer) Shutdown(ctx context.Context) {
	i.logger.Info(ctx).Msg("shutting down indexer")
}
