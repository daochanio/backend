package http

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
)

func (h *httpServer) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		address := r.Header.Get("X-Address")
		token := strings.Split(r.Header.Get("Authorization"), " ")

		if len(token) != 2 || token[0] != "Bearer" {
			h.presentUnathorized(w, r, errors.New("invalid header format"))
			return
		}

		err := h.authenticateUseCase.Execute(ctx, &usecases.AuthenticateInput{
			Address: address,
			SigHex:  token[1],
		})

		if err != nil {
			h.presentUnathorized(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, common.ContextKeyAddress, address)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
