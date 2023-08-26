package distribute

import (
	"context"
	"errors"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/distributor/settings"
	"github.com/daochanio/backend/distributor/usecases"
)

type Distributor interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
}

type distributor struct {
	logger             common.Logger
	settings           settings.Settings
	createDistribution *usecases.Distribute
}

func NewDistributor(logger common.Logger, settings settings.Settings, createDistribution *usecases.Distribute) Distributor {
	return &distributor{
		logger,
		settings,
		createDistribution,
	}
}

func (d *distributor) Start(ctx context.Context) {
	d.logger.Info(ctx).Msg("starting distributor")
	for {
		select {
		case <-ctx.Done():
			d.logger.Info(ctx).Msg("distributor stopped")
			return
		default:
			if err := d.createDistribution.Execute(ctx); err != nil {
				if errors.Is(err, common.ErrNotDistributionTime) {
					time.Sleep(10 * time.Second)
				} else {
					d.logger.Error(ctx).Err(err).Msg("error running distribution")
				}
			}
		}
	}
}

func (d *distributor) Shutdown(ctx context.Context) {
	d.logger.Info(ctx).Msg("shutting down distributor")
}
