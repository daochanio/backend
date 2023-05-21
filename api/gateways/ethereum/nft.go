package ethereum

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	cmn "github.com/daochanio/backend/common"
	"github.com/daochanio/backend/ethereum/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (e *ethereumGateway) GetNFTURI(ctx context.Context, standard string, address string, id string) (string, error) {
	switch strings.ToLower(standard) {
	case "erc721":
		return e.getERC721URI(ctx, common.HexToAddress(address), id)
	case "erc1155":
		return e.getERC1155URI(ctx, common.HexToAddress(address), id)
	default:
		return "", errors.New("invalid nft standard")
	}
}

func (e *ethereumGateway) getERC1155URI(ctx context.Context, address common.Address, tokenId string) (string, error) {
	client, err := ethclient.DialContext(ctx, e.settings.BlockchainURI())

	if err != nil {
		return "", fmt.Errorf("erc1155 uri client %w", err)
	}

	id := new(big.Int)
	id, ok := id.SetString(tokenId, 10)

	if !ok {
		return "", errors.New("erc1155 uri invalid token id")
	}

	instance, err := bindings.NewErc1155(address, client)

	if err != nil {
		return "", fmt.Errorf("erc1155 uri contract %w", err)
	}

	uri, err := cmn.FunctionRetrier(ctx, func() (string, error) {
		uri, err := instance.Uri(&bind.CallOpts{}, id)
		return uri, e.tryWrapRetryable(ctx, "erc1155 uri retry", err)
	})

	if err != nil {
		return "", fmt.Errorf("erc1155 uri %w", err)
	}

	// See https://eips.ethereum.org/EIPS/eip-1155
	// token ids are passed in hexadecimal form and should be padded with zeros until 64 chars long
	hexId := hex.EncodeToString(id.Bytes())
	idLength := len([]rune(hexId))
	for i := 0; i < 64-idLength; i++ {
		hexId = "0" + hexId
	}
	uri = strings.Replace(uri, "{id}", hexId, 1)

	return uri, nil
}

func (e *ethereumGateway) getERC721URI(ctx context.Context, address common.Address, tokenId string) (string, error) {
	client, err := ethclient.DialContext(ctx, e.settings.BlockchainURI())

	if err != nil {
		return "", fmt.Errorf("erc721 uri client %w", err)
	}

	id := new(big.Int)
	id, ok := id.SetString(tokenId, 10)

	if !ok {
		return "", errors.New("erc721 uri invalid token id")
	}

	instance, err := bindings.NewErc721(address, client)

	if err != nil {
		return "", fmt.Errorf("erc721 uri contract %w", err)
	}

	return cmn.FunctionRetrier(ctx, func() (string, error) {
		uri, err := instance.TokenURI(&bind.CallOpts{}, id)
		return uri, e.tryWrapRetryable(ctx, "erc721 uri retry", err)
	})
}
