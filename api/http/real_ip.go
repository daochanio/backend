package http

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/daochanio/backend/common"
)

var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")
var trueClientIP = http.CanonicalHeaderKey("True-Client-IP")

// See https://github.com/go-chi/chi/blob/master/middleware/realip.go
func (h *httpServer) realIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h.logger.Info(r.Context()).Msgf("radr %v", r.RemoteAddr)
		h.logger.Info(r.Context()).Msgf("xff %v", r.Header.Get(xForwardedFor))
		h.logger.Info(r.Context()).Msgf("rip %v", r.Header.Get(xRealIP))
		h.logger.Info(r.Context()).Msgf("tcip %v", r.Header.Get(trueClientIP))

		if ip := getIP(r); ip != "" {
			r.RemoteAddr = ip
		}

		if host, err := parseIP(r); err == nil {
			r.RemoteAddr = host
		}

		ctx := context.WithValue(r.Context(), common.ContextKeyRemoteAddress, r.RemoteAddr)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// if xff header present, pick last addr in comma delimited list
func getIP(r *http.Request) (ip string) {
	xff := r.Header.Get(xForwardedFor)

	if xff == "" {
		return ""
	}

	i := strings.LastIndex(xff, ", ")

	if i+2 >= len(xff) {
		return ""
	}

	if i == -1 {
		return xff
	}

	return xff[i+2:]
}

// validate the ip format
// parse the host out of the ip addr
func parseIP(r *http.Request) (string, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	return host, err
}
