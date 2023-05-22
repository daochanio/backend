package index

import (
	"context"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/settings"
	"github.com/daochanio/backend/indexer/usecases"
)

type Indexer interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
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
			err := i.indexBlocks.Execute(ctx)
			// sleep a bit if we didn't index anything to avoid spamming the blockchain or to recover from transient errors
			if err != nil {
				time.Sleep(time.Duration(i.settings.IntervalSeconds()) * time.Second)
			}
		}
	}
}

// Do any future shutdown resource cleanup here
func (i *indexer) Stop(ctx context.Context) {
	i.logger.Info(ctx).Msg("cleaning up indexer")
}
