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

func (h *httpServer) getCommentsRoute(w http.ResponseWriter, r *http.Request) {
	threadId, err := strconv.ParseInt(chi.URLParam(r, "threadId"), 10, 64)

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	comments, err := h.getCommentsUseCase.Execute(r.Context(), usecases.GetCommentsInput{
		ThreadId: threadId,
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentJSON(w, r, http.StatusOK, toCommentsJson(comments))
}

func (h *httpServer) createCommentRoute(w http.ResponseWriter, r *http.Request) {
	var body commentJson
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	id, err := h.createCommentUseCase.Execute(r.Context(), usecases.CreateCommentInput{
		ParentCommentId: body.ParentCommentId,
		ThreadId:        body.ThreadId,
		Address:         "0xd0147bf60c64b88f3a85425012c129ffdc3e6883", // TODO: get address from auth
		Content:         body.Content,
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

func (h *httpServer) deleteCommentRoute(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "commentId"), 10, 64)

	if err != nil {
		h.presentBadRequest(w, r, err)
	}

	err = h.deleteCommentUseCase.Execute(r.Context(), usecases.DeleteCommentInput{
		Id:             id,
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

func (h *httpServer) createCommentVoteRoute(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "commentId"), 10, 64)
	vote := chi.URLParam(r, "vote")

	if err != nil {
		h.presentBadRequest(w, r, err)
	}

	err = h.createCommentVoteUseCase.Execute(r.Context(), usecases.CreateCommentVoteInput{
		CommentId: id,
		Address:   "0xd0147bf60c64b88f3a85425012c129ffdc3e6883", // TODO: get address from auth
		Vote:      common.VoteType(vote),
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentStatus(w, r, http.StatusCreated)
}

type commentJson struct {
	Id              int64      `json:"id"`
	ParentCommentId *int64     `json:"parentCommentId"`
	ThreadId        int64      `json:"threadId"`
	Address         string     `json:"address"`
	Content         string     `json:"content"`
	IsDeleted       bool       `json:"isDeleted"`
	CreatedAt       time.Time  `json:"createdAt"`
	DeletedAt       *time.Time `json:"deletedAt"`
	Votes           int64      `json:"votes"`
}

func toCommentJson(comment entities.Comment) commentJson {
	return commentJson{
		Id:              comment.GetId(),
		ParentCommentId: comment.GetParentCommentId(),
		ThreadId:        comment.GetThreadId(),
		Address:         comment.GetAddress(),
		Content:         comment.GetContent(),
		IsDeleted:       comment.GetIsDeleted(),
		CreatedAt:       comment.GetCreatedAt(),
		DeletedAt:       comment.GetDeletedAt(),
		Votes:           comment.GetVotes(),
	}
}

func toCommentsJson(comments []entities.Comment) []commentJson {
	commentsJson := make([]commentJson, len(comments))

	for i, comment := range comments {
		commentsJson[i] = toCommentJson(comment)
	}

	return commentsJson
}
