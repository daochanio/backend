package usecases

import (
	"context"
	"fmt"

	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/common"
	"github.com/golang-jwt/jwt/v5"
)

type Authenticate struct {
	settings settings.Settings
}

func NewAuthenticateUseCase(settings settings.Settings) *Authenticate {
	return &Authenticate{
		settings,
	}
}

type AuthenticateInput struct {
	Token string
}

func (u *Authenticate) Execute(ctx context.Context, input *AuthenticateInput) (string, error) {
	if err := common.ValidateStruct(input); err != nil {
		return "", err
	}

	token, err := jwt.Parse(input.Token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(u.settings.JWTSecret()), nil
	})

	if err != nil {
		return "", fmt.Errorf("error parsing token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	return claims.GetSubject()
}
