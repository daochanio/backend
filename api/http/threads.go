package http

import (
	"encoding/json"
	"errors"
	"fmt"
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
		return
	}

	page, err := h.getPage(r)

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	thread, count, err := h.getThreadUseCase.Execute(r.Context(), usecases.GetThreadInput{
		ThreadId:      id,
		CommentOffset: page.Offset,
		CommentLimit:  page.Limit,
	})

	if errors.Is(err, common.ErrNotFound) {
		h.presentNotFound(w, r, err)
		return
	}

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	page.Count = count

	h.presentJSON(w, r, http.StatusOK, toThreadJson(thread), &page)
}

func (h *httpServer) getThreadsRoute(w http.ResponseWriter, r *http.Request) {
	page, err := h.getPage(r)

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	threads, err := h.getThreadsUseCase.Execute(r.Context(), usecases.GetThreadsInput{
		Limit: page.Limit,
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentJSON(w, r, http.StatusOK, toThreadsJson(threads), nil)
}

func (h *httpServer) createThreadRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body createThreadJson
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	thread, err := h.createThreadUseCase.Execute(ctx, usecases.CreateThreadInput{
		Address:       ctx.Value(common.ContextKeyAddress).(string),
		Title:         body.Title,
		ImageFileName: body.ImageFileName,
		Content:       body.Content,
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentJSON(w, r, http.StatusCreated, toThreadJson(thread), nil)
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
	value := chi.URLParam(r, "value")

	if err != nil {
		h.presentBadRequest(w, r, err)
	}

	err = h.createVoteUseCase.Execute(ctx, usecases.CreateVoteInput{
		Id:      id,
		Address: ctx.Value(common.ContextKeyAddress).(string),
		Value:   common.VoteValue(value),
		Type:    common.ThreadVote,
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentStatus(w, r, http.StatusCreated)
}

type createThreadJson struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	ImageFileName string `json:"imageFileName"`
}

type threadJson struct {
	Id        string         `json:"id"`
	Address   string         `json:"address"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Image     *imageJson     `json:"image,omitempty"` // empty if thread deleted
	Comments  *[]commentJson `json:"comments,omitempty"`
	IsDeleted bool           `json:"isDeleted"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt *time.Time     `json:"deletedAt,omitempty"`
	Votes     int64          `json:"votes"`
}

func toThreadJson(thread entities.Thread) threadJson {
	json := threadJson{
		Id:        fmt.Sprint(thread.Id()),
		Address:   thread.Address(),
		Title:     thread.Title(),
		Content:   thread.Content(),
		Image:     toImageJson(thread.Image()),
		IsDeleted: thread.IsDeleted(),
		CreatedAt: thread.CreatedAt(),
		DeletedAt: thread.DeletedAt(),
		Votes:     thread.Votes(),
	}

	if thread.Comments() != nil {
		commentsJson := toCommentsJson(*thread.Comments())
		json.Comments = &commentsJson
	}
	return json
}

func toThreadsJson(threads []entities.Thread) []threadJson {
	json := []threadJson{}
	for _, thread := range threads {
		json = append(json, toThreadJson(thread))
	}
	return json
}
