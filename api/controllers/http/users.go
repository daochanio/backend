package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
	"github.com/go-chi/chi/v5"
)

func (h *httpServer) getUserByAddressRoute(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUser.Execute(r.Context(), usecases.GetUserInput{
		Address: chi.URLParam(r, "address"),
	})

	if errors.Is(err, common.ErrNotFound) {
		h.presentNotFound(w, r, err)
		return
	}

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentJSON(w, r, http.StatusOK, toUserJson(user), nil)
}

type userJson struct {
	Address    string     `json:"address"`
	EnsName    *string    `json:"ensName,omitempty"`
	EnsAvatar  *imageJson `json:"ensAvatar,omitempty"`
	Reputation string     `json:"reputation"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

func toUserJson(user entities.User) userJson {
	return userJson{
		Address:    user.Address(),
		EnsName:    user.EnsName(),
		EnsAvatar:  toImageJson(user.EnsAvatar()),
		Reputation: user.Reputation().String(),
		CreatedAt:  user.CreatedAt(),
		UpdatedAt:  user.UpdatedAt(),
	}
}
