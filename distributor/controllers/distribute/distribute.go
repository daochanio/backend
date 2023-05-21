package distribute

import (
	"context"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/distributor/settings"
)

type Distributor interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
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

func (d *distributor) Start(ctx context.Context) {
	d.logger.Info(ctx).Msg("starting distributor")
	for {
		select {
		case <-ctx.Done():
			d.logger.Info(ctx).Msg("distributor stopped")
			return
		default:
			d.distribute(ctx)
		}
	}
}

func (d *distributor) Stop(ctx context.Context) {
	d.logger.Info(ctx).Msg("cleaning up distributor")
}

func (d *distributor) distribute(ctx context.Context) {
	// TODO: wait until it is past X time UTC and it has been at least Y hours since last run
	time.Sleep(time.Second * 10)

	// TODO: run distribution logic here
}
