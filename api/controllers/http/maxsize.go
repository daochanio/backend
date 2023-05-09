package http

import (
	"net/http"
)

// maxSize is a middleware that caps the size of the request body to the maximum KB specified by maxKB.
func (h *httpServer) maxSize(maxKB int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxKB*1024)
			next.ServeHTTP(w, r)
		})
	}
}
