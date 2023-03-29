package worker

import (
	"context"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/indexer/settings"
)

type Worker struct {
	logger   common.ILogger
	settings settings.ISettings
}

func NewWorker(logger common.ILogger, settings settings.ISettings) *Worker {
	return &Worker{
		logger,
		settings,
	}
}

func (d *Worker) Start(ctx context.Context) error {
	d.logger.Info(ctx).Msg("starting worker")
	for {
		// wait until the next interval
		// i.e if the interval is 5 minutes and the current time is 12:03:45
		// then the next time to run is 12:05:00
		next := time.Now().Truncate(d.settings.Interval()).Add(d.settings.Interval())
		d.logger.Info(ctx).Msgf("waiting until %v", next.Format(time.TimeOnly))
		time.Sleep(time.Until(next))

		// create context from the parent ctx with a timeout of interval
		ctx, cancel := context.WithTimeout(ctx, d.settings.Interval())

		d.logger.Info(ctx).Msg("running indexer")
		err := d.Index(ctx)
		if err != nil {
			d.logger.Error(ctx).Err(err).Msg("indexer failed")
		}
		d.logger.Info(ctx).Msg("indexer completed")

		// cancel the timeout on the context
		cancel()
	}
}

func (d *Worker) Index(ctx context.Context) error {
	return nil
}
