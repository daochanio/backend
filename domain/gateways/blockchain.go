package gateways

import (
	"context"
	"math/big"

	"github.com/daochanio/backend/domain/entities"
)

type BlockchainConfig struct {
	BlockchainURL     string
	ReputationAddress string
}

type Blockchain interface {
	Start(ctx context.Context, config BlockchainConfig)
	Shutdown(ctx context.Context)

	GetNameByAddress(ctx context.Context, address string) (*string, error)
	GetAvatarURIByName(ctx context.Context, name string) (*string, error)
	GetNFTURI(ctx context.Context, standard string, address string, id string) (string, error)

	GetLatestBlockNumber(ctx context.Context) (*big.Int, error)
	GetEvents(context.Context, *big.Int, *big.Int) (entities.Events, error)

	VerifySignature(address string, message string, sigHex string) error
}
