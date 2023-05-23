package http

import (
	"fmt"
	"net/http"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

// ensure the user has an ens name before proceeding
func (h *httpServer) ensName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, ok := ctx.Value(common.ContextKeyUser).(entities.User)

		if !ok {
			h.presentForbidden(w, r, fmt.Errorf("invalid user"))
			return
		}

		if user.EnsName() == nil || *user.EnsName() == "" {
			h.presentForbidden(w, r, fmt.Errorf("ens name required"))
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
