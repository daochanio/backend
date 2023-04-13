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
	id, err := strconv.ParseInt(chi.URLParam(r, "threadId"), 10, 64)

	if err != nil {
		h.presentBadRequest(w, r, err)
	}

	thread, err := h.getThreadUseCase.Execute(r.Context(), usecases.GetThreadInput{
		ThreadId: id,
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

	threads, err := h.getThreadsUseCase.Execute(r.Context(), usecases.GetThreadsInput{
		Limit: 10,
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentJSON(w, r, http.StatusOK, toThreadsJson(threads))
}

func (h *httpServer) createThreadRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body createThreadJson
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	id, err := h.createThreadUseCase.Execute(ctx, usecases.CreateThreadInput{
		Address: ctx.Value(common.ContextKeyAddress).(string),
		Content: body.Content,
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentJSON(w, r, http.StatusCreated, struct {
		Id int64 `json:"id"`
	}{
		Id: id,
	})
}

func (h *httpServer) deleteThreadRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseInt(chi.URLParam(r, "threadId"), 10, 64)

	if err != nil {
		h.presentBadRequest(w, r, err)
	}

	err = h.deleteThreadUseCase.Execute(ctx, usecases.DeleteThreadInput{
		ThreadId:       id,
		DeleterAddress: ctx.Value(common.ContextKeyAddress).(string),
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

func (h *httpServer) createThreadVoteRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseInt(chi.URLParam(r, "threadId"), 10, 64)
	vote := chi.URLParam(r, "vote")

	if err != nil {
		h.presentBadRequest(w, r, err)
	}

	err = h.createThreadVoteUseCase.Execute(ctx, usecases.CreateThreadVoteInput{
		ThreadId: id,
		Address:  ctx.Value(common.ContextKeyAddress).(string),
		Vote:     common.VoteType(vote),
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentStatus(w, r, http.StatusCreated)
}

type createThreadJson struct {
	Content string `json:"content"`
}

type getThreadJson struct {
	Id        int64      `json:"id"`
	Address   string     `json:"address"`
	Content   string     `json:"content"`
	IsDeleted bool       `json:"isDeleted"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	Votes     int64      `json:"votes"`
}

func toThreadJson(thread entities.Thread) getThreadJson {
	return getThreadJson{
		Id:        thread.Id(),
		Address:   thread.Address(),
		Content:   thread.Content(),
		IsDeleted: thread.IsDeleted(),
		CreatedAt: thread.CreatedAt(),
		DeletedAt: thread.DeletedAt(),
		Votes:     thread.Votes(),
	}
}

func toThreadsJson(threads []entities.Thread) []getThreadJson {
	threadsJson := []getThreadJson{}
	for _, thread := range threads {
		threadsJson = append(threadsJson, toThreadJson(thread))
	}
	return threadsJson
}
