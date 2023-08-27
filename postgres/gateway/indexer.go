package gateway

import (
	"context"
	"fmt"
	"math/big"

	"github.com/daochanio/backend/core/entities"
	"github.com/daochanio/backend/postgres/bindings"
	"github.com/jackc/pgx/v5/pgtype"
)

func (g *postgresGateway) GetLastIndexedBlock(ctx context.Context) (*big.Int, error) {
	block, err := g.queries.GetLastIndexedBlock(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting last indexed block: %w", err)
	}

	return block.Int, nil
}

func (g *postgresGateway) UpdateLastIndexedBlock(ctx context.Context, block *big.Int) error {
	return g.queries.UpdateLastIndexedBlock(ctx, pgtype.Numeric{
		Int:   block,
		Valid: true,
	})
}

func (g *postgresGateway) InsertTransferEvents(ctx context.Context, from *big.Int, to *big.Int, transfers []entities.Transfer) error {
	tx, err := g.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := g.queries.WithTx(tx)

	if err := qtx.DeleteTransfers(ctx, bindings.DeleteTransfersParams{
		BlockNumber: pgtype.Numeric{
			Int:   from,
			Valid: true,
		},
		BlockNumber_2: pgtype.Numeric{
			Int:   to,
			Valid: true,
		},
	}); err != nil {
		return fmt.Errorf("failed to delete transfers: %w", err)
	}

	params := []bindings.InsertTransfersParams{}
	for _, transfer := range transfers {
		log := transfer.Log()
		params = append(params, bindings.InsertTransfersParams{
			BlockNumber: pgtype.Numeric{
				Int:   log.BlockNumber(),
				Valid: true,
			},
			TransactionID: log.TransactionId(),
			LogIndex:      int64(log.Index()),
			FromAddress:   transfer.FromAddress(),
			ToAddress:     transfer.ToAddress(),
			Amount: pgtype.Numeric{
				Int:   transfer.Amount(),
				Valid: true,
			},
		})
	}

	if _, err := qtx.InsertTransfers(ctx, params); err != nil {
		return fmt.Errorf("failed to insert transfers: %w", err)
	}

	return tx.Commit(ctx)
}

func (g *postgresGateway) UpdateReputation(ctx context.Context, addresses []string) error {
	tx, err := g.db.Begin(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	qtx := g.queries.WithTx(tx)

	if err := qtx.ZeroReputation(ctx, addresses); err != nil {
		return fmt.Errorf("failed to zero reputation: %w", err)
	}

	if err := qtx.AddReputation(ctx, addresses); err != nil {
		return fmt.Errorf("failed to add reputation: %w", err)
	}

	if err := qtx.DeductReputation(ctx, addresses); err != nil {
		return fmt.Errorf("failed to deduct reputation: %w", err)
	}

	return tx.Commit(ctx)
}
