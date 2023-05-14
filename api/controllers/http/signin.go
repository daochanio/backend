package http

import (
	"encoding/json"
	"net/http"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/usecases"
)

func (h *httpServer) signinRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body signinJsonRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	challenge, err := h.signin.Execute(ctx, usecases.SigninInput{
		Address: body.Address,
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentJSON(w, r, http.StatusOK, toChallengeJson(challenge), nil)
}

type signinJsonRequest struct {
	Address string `json:"address"`
}

type challengeJsonResponse struct {
	Message string `json:"message"`
	Expires int64  `json:"expires"`
}

func toChallengeJson(challenge entities.Challenge) challengeJsonResponse {
	return challengeJsonResponse{
		Message: challenge.Message(),
		Expires: challenge.Expires().Unix(),
	}
}
