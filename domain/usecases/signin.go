package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/gateways"
	"github.com/golang-jwt/jwt/v5"
)

type Signin struct {
	logger     common.Logger
	validator  common.Validator
	database   gateways.Database
	stream     gateways.Stream
	blockchain gateways.Blockchain
}

type SigninInput struct {
	Address   string `validate:"eth_addr"`
	Signature string `validate:"hexadecimal,min=1"`
	JWTSecret string `validate:"required"`
}

func NewSigninUseCase(
	logger common.Logger,
	validator common.Validator,
	database gateways.Database,
	stream gateways.Stream,
	blockchain gateways.Blockchain,
) *Signin {
	return &Signin{
		logger,
		validator,
		database,
		stream,
		blockchain,
	}
}

func (u *Signin) Execute(ctx context.Context, input SigninInput) (string, error) {
	if err := u.validator.ValidateStruct(input); err != nil {
		return "", err
	}

	token, err := u.verifySignature(ctx, input.Address, input.Signature, input.JWTSecret)

	if err != nil {
		return "", fmt.Errorf("invalid signature %w", err)
	}

	err = u.updateUser(ctx, input.Address)

	if err != nil {
		return "", fmt.Errorf("failed to upsert user %w", err)
	}

	return token, err
}

func (u *Signin) verifySignature(ctx context.Context, address string, signature string, jwtSecret string) (string, error) {
	challenge, err := u.database.GetChallengeByAddress(ctx, address)

	if err != nil {
		return "", err
	}

	if err := u.blockchain.VerifySignature(challenge.Address(), challenge.Message(), signature); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "api.daochan.io",
		"sub": address,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	return token.SignedString([]byte(jwtSecret))
}

func (u *Signin) updateUser(ctx context.Context, address string) error {
	err := u.database.UpsertUser(ctx, address)

	if err != nil {
		return err
	}

	return u.stream.PublishSignin(ctx, address)
}
