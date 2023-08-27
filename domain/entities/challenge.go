package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

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
}

func NewChallenge(address string, message string) Challenge {
	return Challenge{
		address: address,
		message: message,
	}
}

func GenerateChallenge(address string) Challenge {
	now := time.Now().Format(time.RFC3339)
	nonce := uuid.New().String()
	return Challenge{
		address: address,
		message: fmt.Sprintf(CHALLENGE_MESSAGE, address, now, nonce),
	}
}

func (c *Challenge) Address() string {
	return c.address
}

func (c *Challenge) Message() string {
	return c.message
}
