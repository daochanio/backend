package usecases

import (
	"context"
	"fmt"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/gateways"
)

type MigrateDatabase struct {
	database gateways.Database
}

func NewMigrateDatabase(database gateways.Database, logger common.Logger) *MigrateDatabase {
	return &MigrateDatabase{
		database: database,
	}
}

func (m *MigrateDatabase) Execute(ctx context.Context, config gateways.DatabaseConfig) error {
	if err := m.database.Migrate(ctx, config); err != nil {
		return fmt.Errorf("err migrating database: %w", err)
	}

	return nil
}
