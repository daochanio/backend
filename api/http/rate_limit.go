package http

import (
	"net/http"

	"github.com/daochanio/backend/api/usecases"
)

func (h *httpServer) rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := h.verifyRateLimitUseCase.Execute(ctx, &usecases.VerifyRateLimitInput{
			IpAddress: r.RemoteAddr,
		})

		if err != nil {
			h.presentTooManyRequests(w, r, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
