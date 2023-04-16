package worker

import (
	"context"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/settings"
	"github.com/daochanio/backend/indexer/usecases"
)

type Worker struct {
	logger           common.ILogger
	settings         settings.ISettings
	subscribeUseCase *usecases.IndexBlocksUseCase
}

func NewWorker(logger common.ILogger, settings settings.ISettings, subscribeUseCase *usecases.IndexBlocksUseCase) *Worker {
	return &Worker{
		logger,
		settings,
		subscribeUseCase,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	w.logger.Info(ctx).Msg("starting worker")
	return w.subscribeUseCase.Execute(ctx)
}
