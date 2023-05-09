package distribute

import (
	"context"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/distributor/settings"
)

type Distributor interface {
	Start(ctx context.Context) error
}

type distributor struct {
	logger   common.Logger
	settings settings.Settings
}

func NewDistributor(logger common.Logger, settings settings.Settings) Distributor {
	return &distributor{
		logger,
		settings,
	}
}

func (d *distributor) Start(ctx context.Context) error {
	d.logger.Info(ctx).Msg("starting distributor")
	for {
		// wait until the next interval
		// i.e if the interval is 5 minutes and the current time is 12:03:45
		// then the next time to run is 12:05:00
		next := time.Now().Truncate(d.settings.Interval()).Add(d.settings.Interval())
		d.logger.Info(ctx).Msgf("waiting until %v", next.Format(time.TimeOnly))
		time.Sleep(time.Until(next))

		// create context from the parent ctx with a timeout of interval
		ctx, cancel := context.WithTimeout(ctx, d.settings.Interval())

		d.logger.Info(ctx).Msg("running distribution")
		err := d.distribute(ctx)
		if err != nil {
			d.logger.Error(ctx).Err(err).Msg("distribution failed")
		}
		d.logger.Info(ctx).Msg("distribution completed")

		// cancel the timeout on the context
		cancel()
	}
}

func (d *distributor) distribute(ctx context.Context) error {
	return nil
}
