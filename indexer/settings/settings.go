package settings

import (
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Settings interface {
	PostgresConfig() *pgxpool.Config
	BlockchainURI() string
	ReputationAddress() string
	ReorgOffset() int64
	IntervalSeconds() int64
	MaxBlockRange() int64
}

type settings struct {
	pgConnectionString string
	blockchainURI      string
	reputationAddress  string
	reorgOffset        int64
	intervalSeconds    int64
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

	maxBlockRange, err := strconv.Atoi(os.Getenv("MAX_BLOCK_RANGE"))
	if err != nil {
		panic(err)
	}

	return &settings{
		pgConnectionString: os.Getenv("PG_CONNECTION_STRING"),
		blockchainURI:      os.Getenv("BLOCKCHAIN_URI"),
		reputationAddress:  os.Getenv("REPUTATION_ADDRESS"),
		reorgOffset:        int64(reorgOffset),
		intervalSeconds:    int64(intervalSeconds),
		maxBlockRange:      int64(maxBlockRange),
	}
}

func (s *settings) PostgresConfig() *pgxpool.Config {
	config, err := pgxpool.ParseConfig(s.pgConnectionString)

	if err != nil {
		panic(err)
	}
	return config
}

func (s *settings) BlockchainURI() string {
	return s.blockchainURI
}

func (s *settings) ReputationAddress() string {
	return s.reputationAddress
}

// the number of blocks to offset by to be resilient to reorgs
func (s *settings) ReorgOffset() int64 {
	return s.reorgOffset
}

func (s *settings) IntervalSeconds() int64 {
	return s.intervalSeconds
}

func (s *settings) MaxBlockRange() int64 {
	return s.maxBlockRange
}
