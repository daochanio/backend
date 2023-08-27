package distribute

import (
	"context"
	"errors"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/usecases"
)

type Distributor interface {
	Start(ctx context.Context, config DistributorConfig)
	Shutdown(ctx context.Context)
}

type distributor struct {
	logger             common.Logger
	createDistribution *usecases.Distribute
}

type DistributorConfig struct {
	Interval time.Duration
}

func NewDistributor(logger common.Logger, createDistribution *usecases.Distribute) Distributor {
	return &distributor{
		logger,
		createDistribution,
	}
}

func (d *distributor) Start(ctx context.Context, config DistributorConfig) {
	d.logger.Info(ctx).Msg("starting distributor")

	for {
		select {
		case <-ctx.Done():
			d.logger.Info(ctx).Msg("distributor stopped")
			return
		default:
			if err := d.createDistribution.Execute(ctx, usecases.DistributeInput{
				Interval: config.Interval,
			}); err != nil {
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
