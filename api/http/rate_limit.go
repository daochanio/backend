package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/daochanio/backend/api/usecases"
)

func (h *httpServer) rateLimit(namespace string, rate int, period time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			fmt.Println("rateLimit", namespace, rate, period)

			err := h.verifyRateLimitUseCase.Execute(ctx, &usecases.VerifyRateLimitInput{
				Namespace: namespace,
				IpAddress: r.RemoteAddr,
				Rate:      rate,
				Period:    period,
			})

			if err != nil {
				h.presentTooManyRequests(w, r, err)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
