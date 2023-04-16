package worker

import (
	"context"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/settings"
	"github.com/daochanio/backend/indexer/usecases"
)

type Worker struct {
	logger             common.ILogger
	settings           settings.ISettings
	indexBlocksUseCase *usecases.IndexBlocksUseCase
}

func NewWorker(logger common.ILogger, settings settings.ISettings, indexBlocksUseCase *usecases.IndexBlocksUseCase) *Worker {
	return &Worker{
		logger,
		settings,
		indexBlocksUseCase,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	w.logger.Info(ctx).Msg("starting worker")
	for {
		err := w.indexBlocksUseCase.Execute(ctx)

		// sleep a bit if we didn't index anything to avoid spamming the blockchain or to recover from transient errors
		if err != nil {
			time.Sleep(time.Duration(w.settings.IntervalSeconds()) * time.Second)
		}
	}
}
