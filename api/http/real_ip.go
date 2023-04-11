package http

import (
	"context"
	"net"
	"net/http"

	"github.com/daochanio/backend/common"
)

var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")

// See https://github.com/go-chi/chi/blob/master/middleware/realip.go
func (h *httpServer) realIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h.logger.Info(r.Context()).Msgf("radr %v", r.RemoteAddr)
		h.logger.Info(r.Context()).Msgf("xff %v", r.Header.Get(xForwardedFor))
		h.logger.Info(r.Context()).Msgf("rip %v", r.Header.Get(h.settings.RealIPHeader()))

		if ip := r.Header.Get(h.settings.RealIPHeader()); ip != "" {
			r.RemoteAddr = ip
		}

		if host, err := parseIP(r); err == nil {
			r.RemoteAddr = host
		}

		ctx := context.WithValue(r.Context(), common.ContextKeyRemoteAddress, r.RemoteAddr)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validate the ip format
// parse the host out of the ip addr
func parseIP(r *http.Request) (string, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	return host, err
}
