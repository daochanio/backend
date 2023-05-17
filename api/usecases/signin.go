package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/common"
	"github.com/golang-jwt/jwt/v5"
)

type Signin struct {
	logger   common.Logger
	settings settings.Settings
	database Database
	stream   Stream
}

type SigninInput struct {
	Address   string `validate:"eth_addr"`
	Signature string `validate:"hexadecimal,min=1"`
}

func NewSigninUseCase(logger common.Logger, settings settings.Settings, database Database, stream Stream) *Signin {
	return &Signin{
		logger,
		settings,
		database,
		stream,
	}
}

func (u *Signin) Execute(ctx context.Context, input SigninInput) (string, error) {
	if err := common.ValidateStruct(input); err != nil {
		return "", err
	}

	token, err := u.verifySignature(ctx, input.Address, input.Signature)

	if err != nil {
		return "", fmt.Errorf("invalid signature %w", err)
	}

	err = u.updateUser(ctx, input.Address)

	if err != nil {
		return "", fmt.Errorf("failed to upsert user %w", err)
	}

	return token, err
}

func (u *Signin) verifySignature(ctx context.Context, address string, signature string) (string, error) {
	challenge, err := u.database.GetChallengeByAddress(ctx, address)

	if err != nil {
		return "", err
	}

	if err := challenge.Verify(signature); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "api.daochan.io",
		"sub": address,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	return token.SignedString([]byte(u.settings.JWTSecret()))
}

func (u *Signin) updateUser(ctx context.Context, address string) error {
	err := u.database.UpsertUser(ctx, address)

	if err != nil {
		return err
	}

	return u.stream.PublishSignin(ctx, address)
}
