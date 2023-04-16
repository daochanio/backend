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
		// TODO: If this takes longer than a minute we are at risk of falling behind the latest state
		// We should probably alert if we notice this happening
		w.indexBlocksUseCase.Execute(ctx)
		time.Sleep(1 * time.Minute)
	}
}
