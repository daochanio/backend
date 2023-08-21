package usecases

import (
	"context"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/distributor/settings"
)

type CreateDistribution struct {
	settings settings.Settings
	logger   common.Logger
}

func NewCreateDistribution(
	settings settings.Settings,
	logger common.Logger,
) *CreateDistribution {
	return &CreateDistribution{
		settings,
		logger,
	}
}

// Read the last distribution record from the database or create one if it doesn't exist.
// Calculate the next earliest time to run a new distribution based on the interval.
// If the last run time is within the next run interval, we do not run a new distribution.
// If the alst run time is not within the next run interval, we do run a new distribution.
// E.g Suppose, the run interval is 5 minutes.
// If the last run distribution time was recorded at 12:03:23 then the next time to run is any time greater than 12:05:00.
// If the current time is 12:03:45, then we should not run.
// But if the current time is 12:05:07, then we should run.
//
// When the next distribution round runs, the votes that are accepted and not associated
// with a distribution are processed and then tied to a distribution through FK.
func (u *CreateDistribution) Execute(ctx context.Context) error {
	next := time.Now().Truncate(u.settings.Interval()).Add(u.settings.Interval())

	u.logger.Info(ctx).Msgf("next distribution will run in %s", time.Until(next))

	return common.ErrNotDistributionTime
}
