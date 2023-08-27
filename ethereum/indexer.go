package ethereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/daochanio/backend/core/entities"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func (g *ethereumGateway) GetLatestBlockNumber(ctx context.Context) (*big.Int, error) {
	header, err := g.ethClient.HeaderByNumber(ctx, nil)
	if err == ethereum.NotFound {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get latest block number: %w", err)
	}

	return header.Number, nil
}

func (g *ethereumGateway) GetEvents(ctx context.Context, fromBlock *big.Int, toBlock *big.Int) (entities.Events, error) {
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: g.eventSources,
		Topics: [][]common.Hash{
			{
				crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")),
			},
		},
	}

	logs, err := g.ethClient.FilterLogs(ctx, query)

	if err != nil {
		return entities.Events{}, err
	}

	transferEvents := []entities.Transfer{}
	for _, log := range logs {
		g.logger.Info(ctx).Msgf("found log for address: %v at block: %d at index %d", log.Address.Hex(), log.BlockNumber, log.Index)

		if transfer, err := g.toTransfer(log); err == nil {
			transferEvents = append(transferEvents, transfer)
		} else {
			return entities.Events{}, fmt.Errorf("failed to parse log into any event: %w", err)
		}
	}

	return entities.NewEvents(transferEvents), nil
}

func (g *ethereumGateway) toTransfer(log types.Log) (entities.Transfer, error) {
	transfer, err := g.reputation.ParseTransfer(log)
	if err != nil {
		return entities.Transfer{}, fmt.Errorf("failed to parse event into transfer: %w", err)
	}

	if transfer.Value == nil {
		return entities.Transfer{}, errors.New("nil value in transfer log")
	}
	return entities.NewTransfer(
		transfer.From.Hex(),
		transfer.To.Hex(),
		transfer.Value,
		g.toLog(log),
	), nil
}

func (g *ethereumGateway) toLog(log types.Log) entities.Log {
	return entities.NewLog(
		new(big.Int).SetUint64(log.BlockNumber),
		log.TxHash.Hex(),
		uint32(log.Index),
	)
}
