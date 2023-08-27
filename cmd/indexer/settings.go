package main

import (
	"os"
	"strconv"
	"time"

	"github.com/daochanio/backend/cmd/indexer/index"
	"github.com/daochanio/backend/domain/gateways"
)

type Settings interface {
	IndexerConfig() index.IndexerConfig
	DatabaseConfig() gateways.DatabaseConfig
	BlockchainConfig() gateways.BlockchainConfig
}

type settings struct {
	pgConnectionString string
	blockchainURL      string
	reputationAddress  string
	reorgOffset        int64
	interval           time.Duration
	maxBlockRange      int64
}

func NewSettings() Settings {
	reorgOffset, err := strconv.Atoi(os.Getenv("REORG_OFFSET"))
	if err != nil {
		panic(err)
	}

	intervalSeconds, err := strconv.Atoi(os.Getenv("INTERVAL_SECONDS"))
	if err != nil {
		panic(err)
	}
	interval := time.Duration(intervalSeconds) * time.Second

	maxBlockRange, err := strconv.Atoi(os.Getenv("MAX_BLOCK_RANGE"))
	if err != nil {
		panic(err)
	}

	return &settings{
		pgConnectionString: os.Getenv("PG_CONNECTION_STRING"),
		blockchainURL:      os.Getenv("BLOCKCHAIN_URI"),
		reputationAddress:  os.Getenv("REPUTATION_ADDRESS"),
		reorgOffset:        int64(reorgOffset),
		interval:           interval,
		maxBlockRange:      int64(maxBlockRange),
	}
}

func (s *settings) IndexerConfig() index.IndexerConfig {
	return index.IndexerConfig{
		Interval:      s.interval,
		MaxBlockRange: s.maxBlockRange,
		ReorgOffset:   s.reorgOffset,
	}
}

func (s *settings) DatabaseConfig() gateways.DatabaseConfig {
	return gateways.DatabaseConfig{
		ConnectionString: s.pgConnectionString,
		MinConnections:   10,
		MaxConnections:   100,
	}
}

func (s *settings) BlockchainConfig() gateways.BlockchainConfig {
	return gateways.BlockchainConfig{
		BlockchainURL: s.blockchainURL,
	}
}
