package settings

import (
	"os"
	"strconv"
)

type ISettings interface {
	DbConnectionString() string
	BlockchainURI() string
	TokenAddress() string
	GovernorAddress() string
	ControllerAddress() string
	ReorgOffset() int64
	IntervalSeconds() int64
}

type settings struct {
	dbConnectionString string
	blockchainURI      string
	tokenAddress       string
	governorAddress    string
	controllerAddress  string
	reorgOffset        int64
	intervalSeconds    int64
}

func NewSettings() ISettings {
	reorgOffset, err := strconv.Atoi(os.Getenv("REORG_OFFSET"))
	if err != nil {
		panic(err)
	}

	intervalSeconds, err := strconv.Atoi(os.Getenv("INTERVAL_SECONDS"))
	if err != nil {
		panic(err)
	}

	return &settings{
		dbConnectionString: os.Getenv("PG_CONNECTION_STRING"),
		blockchainURI:      os.Getenv("BLOCKCHAIN_URI"),
		tokenAddress:       os.Getenv("TOKEN_ADDRESS"),
		governorAddress:    os.Getenv("GOVERNOR_ADDRESS"),
		controllerAddress:  os.Getenv("CONTROLLER_ADDRESS"),
		reorgOffset:        int64(reorgOffset),
		intervalSeconds:    int64(intervalSeconds),
	}
}

func (s *settings) DbConnectionString() string {
	return s.dbConnectionString
}

func (s *settings) BlockchainURI() string {
	return s.blockchainURI
}

func (s *settings) TokenAddress() string {
	return s.tokenAddress
}

func (s *settings) GovernorAddress() string {
	return s.governorAddress
}

func (s *settings) ControllerAddress() string {
	return s.controllerAddress
}

// the number of blocks to offset by to be resilient to reorgs
func (s *settings) ReorgOffset() int64 {
	return s.reorgOffset
}

func (s *settings) IntervalSeconds() int64 {
	return s.intervalSeconds
}
