package http

import (
	"context"
	"net/http"
	"time"

	"github.com/daochanio/backend/common"
)

// Record the start time of the request
// ContextKeyRequestStartTime will be consumed by the http presenter when logging the request event
func (h *httpServer) timer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		ctx := context.WithValue(r.Context(), common.ContextKeyRequestStartTime, t1)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
