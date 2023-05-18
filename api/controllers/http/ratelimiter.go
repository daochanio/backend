package http

import (
	"net/http"
	"time"

	"github.com/daochanio/backend/api/usecases"
)

func (h *httpServer) rateLimiter(namespace string, rate int, period time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			err := h.rateLimit.Execute(ctx, &usecases.RateLimitInput{
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
