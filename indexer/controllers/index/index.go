package index

import (
	"context"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/settings"
	"github.com/daochanio/backend/indexer/usecases"
)

type Indexer interface {
	Start(ctx context.Context) error
}

type indexer struct {
	logger             common.Logger
	settings           settings.Settings
	indexBlocksUseCase *usecases.IndexBlocksUseCase
}

func NewIndexer(logger common.Logger, settings settings.Settings, indexBlocksUseCase *usecases.IndexBlocksUseCase) Indexer {
	return &indexer{
		logger,
		settings,
		indexBlocksUseCase,
	}
}

func (i *indexer) Start(ctx context.Context) error {
	i.logger.Info(ctx).Msg("starting worker")
	for {
		err := i.indexBlocksUseCase.Execute(ctx)

		// sleep a bit if we didn't index anything to avoid spamming the blockchain or to recover from transient errors
		if err != nil {
			time.Sleep(time.Duration(i.settings.IntervalSeconds()) * time.Second)
		}
	}
}
