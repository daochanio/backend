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
	cache    Cache
}

func NewAuthenticateUseCase(settings settings.Settings, cache Cache) *Authenticate {
	return &Authenticate{
		settings,
		cache,
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

	claims := token.Claims.(jwt.MapClaims)

	return claims.GetSubject()
}
