package http

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/daochanio/backend/common"
)

// add a randomly generated trace id to the context
//
// see https://github.com/go-chi/chi/blob/master/middleware/request_id.go#L67
func (h *httpServer) traceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := fmt.Sprintf("%10d", rand.Intn(10000000000))
		ctx = context.WithValue(ctx, common.ContextKeyTraceID, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
