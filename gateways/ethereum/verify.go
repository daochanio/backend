package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Ref: https://gist.github.com/dcb9/385631846097e1f59e3cba3b1d42f3ed#file-eth_sign_verify-go
func (g *ethereumGateway) VerifySignature(address string, message string, sigHex string) error {
	msg := []byte(message)
	msg = accounts.TextHash(msg)
	sig, err := hexutil.Decode(sigHex)

	if err != nil {
		return fmt.Errorf("decoding signature: %w", err)
	}

	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}

	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return fmt.Errorf("recovering public key from signature: %w", err)
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	if address != recoveredAddr.Hex() {
		return fmt.Errorf("signature does not match challenge address")
	}

	return nil
}
