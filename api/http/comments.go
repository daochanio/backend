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

func (h *httpServer) getCommentsRoute(w http.ResponseWriter, r *http.Request) {
	threadId, err := strconv.ParseInt(chi.URLParam(r, "threadId"), 10, 64)

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	page, err := h.getPage(r)

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	comments, count, err := h.getCommentsUseCase.Execute(r.Context(), usecases.GetCommentsInput{
		ThreadId: threadId,
		Offset:   page.Offset,
		Limit:    page.Limit,
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	page.Count = count

	h.presentJSON(w, r, http.StatusOK, toCommentsJson(comments), &page)
}

func (h *httpServer) createCommentRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body createCommentJson
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.presentBadRequest(w, r, fmt.Errorf("invalid json: %w", err))
		return
	}

	var repliedToCommentId *int64
	if body.RepliedToCommentId != nil {
		id, err := strconv.ParseInt(*body.RepliedToCommentId, 10, 64)

		if err != nil {
			h.presentBadRequest(w, r, fmt.Errorf("invalid replied to comment id: %w", err))
			return
		}

		repliedToCommentId = &id
	}

	threadId, err := strconv.ParseInt(body.ThreadId, 10, 64)

	if err != nil {
		h.presentBadRequest(w, r, fmt.Errorf("invalid thread id: %w", err))
		return
	}

	comment, err := h.createCommentUseCase.Execute(ctx, usecases.CreateCommentInput{
		RepliedToCommentId: repliedToCommentId,
		ThreadId:           threadId,
		Address:            ctx.Value(common.ContextKeyAddress).(string),
		Content:            body.Content,
		ImageFileName:      body.ImageFileName,
	})

	if err != nil {
		h.presentBadRequest(w, r, fmt.Errorf("invalid comment: %w", err))
		return
	}

	h.presentJSON(w, r, http.StatusCreated, toCommentJson(comment), nil)
}

func (h *httpServer) deleteCommentRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(chi.URLParam(r, "commentId"), 10, 64)

	if err != nil {
		h.presentBadRequest(w, r, err)
	}

	err = h.deleteCommentUseCase.Execute(ctx, usecases.DeleteCommentInput{
		Id:             id,
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

func (h *httpServer) createCommentVoteRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(chi.URLParam(r, "commentId"), 10, 64)
	value := chi.URLParam(r, "value")

	if err != nil {
		h.presentBadRequest(w, r, err)
	}

	err = h.createVoteUseCase.Execute(ctx, usecases.CreateVoteInput{
		Id:      id,
		Address: ctx.Value(common.ContextKeyAddress).(string),
		Value:   common.VoteValue(value),
		Type:    common.CommentVote,
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentStatus(w, r, http.StatusCreated)
}

type createCommentJson struct {
	RepliedToCommentId *string `json:"repliedToCommentId,omitempty"`
	ThreadId           string  `json:"threadId"`
	Content            string  `json:"content"`
	ImageFileName      string  `json:"imageFileName"`
}

type commentJson struct {
	Id               string       `json:"id"`
	RepliedToComment *commentJson `json:"repliedToComment,omitempty"`
	ThreadId         string       `json:"threadId,omitempty"` // empty if reply
	Address          string       `json:"address"`
	Content          string       `json:"content"`
	Image            *imageJson   `json:"image,omitempty"` // empty if comment deleted
	IsDeleted        bool         `json:"isDeleted"`
	CreatedAt        time.Time    `json:"createdAt"`
	DeletedAt        *time.Time   `json:"deletedAt,omitempty"` // empty if comment not deleted
	Votes            int64        `json:"votes,omitempty"`     // empty if reply
}

func toCommentsJson(comments []entities.Comment) []commentJson {
	json := make([]commentJson, len(comments))

	for i, comment := range comments {
		json[i] = toCommentJson(comment)
	}

	return json
}

func toCommentJson(comment entities.Comment) commentJson {
	json := commentJson{
		Id:        fmt.Sprint(comment.Id()),
		ThreadId:  fmt.Sprint(comment.ThreadId()),
		Address:   comment.Address(),
		Content:   comment.Content(),
		Image:     toImageJson(comment.Image()),
		IsDeleted: comment.IsDeleted(),
		CreatedAt: comment.CreatedAt(),
		DeletedAt: comment.DeletedAt(),
		Votes:     comment.Votes(),
	}

	if repliedToComment := comment.RepliedToComment(); repliedToComment != nil {
		repliedToCommentJson := commentJson{
			Id:        fmt.Sprint(repliedToComment.Id()),
			Address:   repliedToComment.Address(),
			Content:   repliedToComment.Content(),
			Image:     toImageJson(repliedToComment.Image()),
			IsDeleted: repliedToComment.IsDeleted(),
			CreatedAt: repliedToComment.CreatedAt(),
			DeletedAt: repliedToComment.DeletedAt(),
		}
		json.RepliedToComment = &repliedToCommentJson
	}

	return json
}
