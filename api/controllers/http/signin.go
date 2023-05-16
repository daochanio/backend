package http

import (
	"encoding/json"
	"net/http"

	"github.com/daochanio/backend/api/usecases"
	"github.com/go-chi/chi/v5"
)

func (h *httpServer) getChallengeRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	address := chi.URLParam(r, "address")

	challenge, err := h.getChallenge.Execute(ctx, &usecases.GetChallengeInput{
		Address: address,
	})

	if err != nil {
		h.presentUnathorized(w, r, err)
		return
	}

	h.presentJSON(w, r, http.StatusOK, challengeJsonResponse{
		Message: challenge.Message(),
	}, nil)
}

func (h *httpServer) putChallengeRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	address := chi.URLParam(r, "address")

	var body signinJsonRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.presentUnathorized(w, r, err)
		return
	}

	token, err := h.signin.Execute(ctx, usecases.SigninInput{
		Address:   address,
		Signature: body.Signature,
	})

	if err != nil {
		h.presentUnathorized(w, r, err)
		return
	}

	h.presentJSON(w, r, http.StatusOK, signinJsonResponse{
		Token: token,
	}, nil)
}

type challengeJsonResponse struct {
	Message string `json:"message"`
}

type signinJsonRequest struct {
	Signature string `json:"signature"`
}

type signinJsonResponse struct {
	Token string `json:"token"`
}
