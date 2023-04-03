package entities

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// how long a signature will ultimately be valid for
const CHALLENGE_TTL = time.Hour * 24

// the message that will be presented to the user in their wallet when signing
const CHALLENGE_MESSAGE string = `
Please sign this message to prove you are the owner of this address: %v

Signing this message does not cost anything.

Timestamp: %v
Nonce: %v
`

type Challenge struct {
	address string
	message string
	ttl     time.Duration
	expires time.Time
}

func NewChallenge(address string, message string, expires time.Time) Challenge {
	return Challenge{
		address: address,
		message: message,
		expires: expires,
	}
}

func GenerateChallenge(address string) Challenge {
	now := time.Now()
	ttl := CHALLENGE_TTL
	nonce := fmt.Sprintf("%10d", rand.Intn(10000000000))
	return Challenge{
		address: address,
		message: fmt.Sprintf(CHALLENGE_MESSAGE, address, now.Format(time.RFC3339), nonce),
		ttl:     ttl,
		expires: now.Add(ttl),
	}
}

func (c *Challenge) Address() string {
	return c.address
}

func (c *Challenge) Message() string {
	return c.message
}

func (c *Challenge) TTL() time.Duration {
	return c.ttl
}

func (c *Challenge) Expires() time.Time {
	return c.expires
}

// Shoutout: https://gist.github.com/dcb9/385631846097e1f59e3cba3b1d42f3ed#file-eth_sign_verify-go
func (c *Challenge) Verify(sigHex string) error {
	msg := []byte(c.message)
	msg = accounts.TextHash(msg)
	sig := hexutil.MustDecode(sigHex)

	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}

	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return fmt.Errorf("recovering public key from signature: %w", err)
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	if c.address != recoveredAddr.Hex() {
		return fmt.Errorf("signature does not match challenge address")
	}

	return nil
}
