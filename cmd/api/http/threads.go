package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/domain/usecases"
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

	thread, count, err := h.getThread.Execute(r.Context(), usecases.GetThreadInput{
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

	threads, err := h.getThreads.Execute(r.Context(), usecases.GetThreadsInput{
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
	body, err := common.Decode[createThreadJson](r.Body)
	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	user, ok := ctx.Value(common.ContextKeyUser).(entities.User)

	if !ok {
		h.presentBadRequest(w, r, fmt.Errorf("invalid user"))
		return
	}

	thread, err := h.createThread.Execute(ctx, usecases.CreateThreadInput{
		Address:       user.Address(),
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
		return
	}

	user, ok := ctx.Value(common.ContextKeyUser).(entities.User)

	if !ok {
		h.presentBadRequest(w, r, fmt.Errorf("invalid user"))
		return
	}

	err = h.deleteThread.Execute(ctx, usecases.DeleteThreadInput{
		ThreadId:       id,
		DeleterAddress: user.Address(),
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

	user, ok := ctx.Value(common.ContextKeyUser).(entities.User)

	if !ok {
		h.presentBadRequest(w, r, fmt.Errorf("invalid user"))
		return
	}

	err = h.createVote.Execute(ctx, usecases.CreateVoteInput{
		Id:      id,
		Address: user.Address(),
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
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Image     *imageJson     `json:"image,omitempty"` // empty if thread deleted
	User      userJson       `json:"user"`
	Comments  *[]commentJson `json:"comments,omitempty"`
	IsDeleted bool           `json:"isDeleted"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt *time.Time     `json:"deletedAt,omitempty"`
	Votes     int64          `json:"votes"`
}

func toThreadJson(thread entities.Thread) threadJson {
	json := threadJson{
		Id:        fmt.Sprint(thread.Id()),
		Title:     thread.Title(),
		Content:   thread.Content(),
		Image:     toImageJson(thread.Image()),
		User:      toUserJson(thread.User()),
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
