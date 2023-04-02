package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
	"github.com/go-chi/chi/v5"
)

func (h *httpServer) getThreadByIdRoute(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)

	if err != nil {
		h.presentBadRequest(w, r, err)
	}

	thread, err := h.getThreadUseCase.Execute(r.Context(), usecases.GetThreadInput{
		ThreadId: int32(id),
	})

	if errors.Is(err, common.ErrNotFound) {
		h.presentNotFound(w, r, err)
		return
	}

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentJSON(w, r, http.StatusOK, toThreadJson(thread))
}

func (h *httpServer) getThreadsRoute(w http.ResponseWriter, r *http.Request) {
	paginationParams, err := h.getPaginationParams(r)
	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	threads, err := h.getThreadsUseCase.Execute(r.Context(), usecases.GetThreadsInput{
		Offset: paginationParams.Offset,
		Limit:  paginationParams.Limit,
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentJSON(w, r, http.StatusOK, toThreadsJson(threads))
}

func (h *httpServer) createThreadRoute(w http.ResponseWriter, r *http.Request) {
	var input usecases.CreateThreadInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	input.Address = "0xd0147bf60c64b88f3a85425012c129ffdc3e6883" // TODO: get address from auth

	id, err := h.createThreadUseCase.Execute(r.Context(), input)

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentJSON(w, r, http.StatusCreated, struct {
		Id int32 `json:"id"`
	}{
		Id: id,
	})
}

func (h *httpServer) deleteThreadRoute(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)

	if err != nil {
		h.presentBadRequest(w, r, err)
	}

	err = h.deleteThreadUseCase.Execute(r.Context(), usecases.DeleteThreadInput{
		ThreadId:       int32(id),
		DeleterAddress: "0xd0147bf60c64b88f3a85425012c129ffdc3e6883", // TODO: get address from auth
	})

	if errors.Is(err, common.ErrNotFound) {
		h.presentNotFound(w, r, err)
		return
	}

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentStatus(w, r, http.StatusOK)
}

func (h *httpServer) voteThreadRoute(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)

	vote := chi.URLParam(r, "vote")

	if err != nil {
		h.presentBadRequest(w, r, err)
	}

	err = h.voteThreadUseCase.Execute(r.Context(), usecases.VoteThreadInput{
		ThreadId: int32(id),
		Address:  "0xd0147bf60c64b88f3a85425012c129ffdc3e6883", // TODO: get address from auth
		Vote:     usecases.VoteType(vote),
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentStatus(w, r, http.StatusCreated)
}

type threadJson struct {
	Id        int32      `json:"id"`
	Address   string     `json:"address"`
	Content   string     `json:"content"`
	IsDeleted bool       `json:"isDeleted"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	Votes     int64      `json:"votes"`
}

func toThreadJson(thread entities.Thread) threadJson {
	return threadJson{
		Id:        thread.GetId(),
		Address:   thread.GetAddress(),
		Content:   thread.GetContent(),
		IsDeleted: thread.GetIsDeleted(),
		CreatedAt: thread.GetCreatedAt(),
		DeletedAt: thread.GetDeletedAt(),
		Votes:     thread.GetVotes(),
	}
}

func toThreadsJson(threads []entities.Thread) []threadJson {
	threadsJson := []threadJson{}
	for _, thread := range threads {
		threadsJson = append(threadsJson, toThreadJson(thread))
	}
	return threadsJson
}
