package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/daochanio/backend/api/settings"
	"github.com/daochanio/backend/api/usecases"
	"github.com/daochanio/backend/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type HttpServer interface {
	Start(context.Context) error
}

type httpServer struct {
	logger                 common.Logger
	settings               settings.Settings
	signinUseCase          *usecases.SigninUseCase
	verifyChallengeUseCase *usecases.VerifyChallengeUseCase
	verifyRateLimitUseCase *usecases.VerifyRateLimitUseCase
	createThreadUseCase    *usecases.CreateThreadUseCase
	getThreadUseCase       *usecases.GetThreadUseCase
	getThreadsUseCase      *usecases.GetThreadsUseCase
	deleteThreadUseCase    *usecases.DeleteThreadUseCase
	createVoteUseCase      *usecases.CreateVoteUseCase
	createCommentUseCase   *usecases.CreateCommentUseCase
	getCommentsUseCase     *usecases.GetCommentsUseCase
	deleteCommentUseCase   *usecases.DeleteCommentUseCase
	uploadImageUseCase     *usecases.UploadImageUsecase
}

func NewHttpServer(
	logger common.Logger,
	settings settings.Settings,
	signinUseCase *usecases.SigninUseCase,
	verifyChallengeUseCase *usecases.VerifyChallengeUseCase,
	verifyRateLimitUseCase *usecases.VerifyRateLimitUseCase,
	createThreadUseCase *usecases.CreateThreadUseCase,
	getThreadUseCase *usecases.GetThreadUseCase,
	getThreadsUseCase *usecases.GetThreadsUseCase,
	deleteThreadUseCase *usecases.DeleteThreadUseCase,
	createVoteUseCase *usecases.CreateVoteUseCase,
	createCommentUseCase *usecases.CreateCommentUseCase,
	getCommentsUseCase *usecases.GetCommentsUseCase,
	deleteCommentUseCase *usecases.DeleteCommentUseCase,
	uploadImageUseCase *usecases.UploadImageUsecase) HttpServer {
	return &httpServer{
		logger,
		settings,
		signinUseCase,
		verifyChallengeUseCase,
		verifyRateLimitUseCase,
		createThreadUseCase,
		getThreadUseCase,
		getThreadsUseCase,
		deleteThreadUseCase,
		createVoteUseCase,
		createCommentUseCase,
		getCommentsUseCase,
		deleteCommentUseCase,
		uploadImageUseCase,
	}
}

func (h *httpServer) Start(ctx context.Context) error {
	h.logger.Info(ctx).Msg("starting http service")

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://daochan.io", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Address"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.NoCache)
	r.Use(h.timer)
	r.Use(h.realIP)
	r.Use(h.traceID)
	r.Use(h.recoverer)
	r.Use(h.timeout)

	r.Get("/", h.healthRoute)

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.Compress(5, "application/json"))

		// public routes
		r.Group(func(r chi.Router) {
			r.Use(h.rateLimit("public", 10000, time.Minute)) // TODO:  Make rate limiting more restrictive
			r.Use(h.maxSize(1))

			r.Put("/signin", h.signinRoute)
			r.Get("/threads", h.getThreadsRoute)
			r.Get("/threads/{threadId}", h.getThreadByIdRoute)
			r.Get("/threads/{threadId}/comments", h.getCommentsRoute)
		})

		// authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(h.authenticate)
			r.Use(h.maxSize(5))

			r.With(h.rateLimit("create:thread", 2, time.Minute*10)).Post("/threads", h.createThreadRoute)
			r.With(h.rateLimit("vote:thread", 10000, time.Minute)).Put("/threads/{threadId}/votes/{value}", h.createThreadVoteRoute) // TODO:  Make rate limiting more restrictive
			r.With(h.rateLimit("create:comment", 5, time.Minute*10)).Post("/threads/{threadId}/comments", h.createCommentRoute)
			r.With(h.rateLimit("vote:comment", 10, time.Minute)).Put("/threads/{threadId}/comments/{commentId}/votes/{value}", h.createCommentVoteRoute)
		})

		// permissioned routes
		r.Group(func(r chi.Router) {
			r.Use(h.authenticate) // TODO: Make this require moderator+ permission
			r.Use(h.rateLimit("permissioned", 10, time.Second))
			r.Use(h.maxSize(1))

			r.Delete("/threads/{threadId}", h.deleteThreadRoute)
			r.Delete("/threads/{threadId}/comments/{commentId}", h.deleteCommentRoute)
		})

		// image route
		r.Group(func(r chi.Router) {
			r.Use(h.authenticate)
			r.Use(h.rateLimit("create:image", 7, time.Minute*10)) // should encompass creating images for threads and comments
			r.Use(h.maxSize(3 * 1024))

			r.Post("/images", h.uploadImageRoute)
		})
	})

	port := h.settings.Port()

	h.logger.Info(ctx).Msgf("listening on port %v", port)

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), r)

	h.logger.Error(ctx).Err(err).Msg("error in http service")

	return err
}

func (h *httpServer) presentNotFound(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Warn(r.Context()).Err(err).Msg("not found")
	h.presentJSON(w, r, http.StatusNotFound, toErrJson("not found"), nil)
}

func (h *httpServer) presentBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Warn(r.Context()).Err(err).Msg("bad request")
	h.presentJSON(w, r, http.StatusBadRequest, toErrJson("bad request"), nil)
}

func (h *httpServer) presentUnathorized(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Warn(r.Context()).Err(err).Msg("unauthorized")
	h.presentJSON(w, r, http.StatusUnauthorized, toErrJson("unathorized"), nil)
}

func (h *httpServer) presentTooManyRequests(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Warn(r.Context()).Err(err).Msg("too many requests")
	h.presentJSON(w, r, http.StatusTooManyRequests, toErrJson("too many requests"), nil)
}

func (h *httpServer) presentJSON(w http.ResponseWriter, r *http.Request, statusCode int, data any, lastPage *pageJson) {
	w.Header().Set("Content-Type", "application/json")
	h.presentStatus(w, r, statusCode)

	var nextPage *pageJson
	if lastPage != nil && lastPage.Offset+lastPage.Limit < lastPage.Count {
		nextPage = &pageJson{
			Offset: lastPage.Offset + lastPage.Limit,
			Limit:  lastPage.Limit,
			Count:  lastPage.Count,
		}
	}
	if err := json.NewEncoder(w).Encode(bodyJson{
		Data:     data,
		NextPage: nextPage,
	}); err != nil {
		h.logger.Error(r.Context()).Err(err).Msg("error encoding json")
	}
}

func (h *httpServer) presentText(w http.ResponseWriter, r *http.Request, statusCode int, text string) {
	w.Header().Set("Content-Type", "text/plain")
	h.presentStatus(w, r, statusCode)
	if _, err := w.Write([]byte(text)); err != nil {
		h.logger.Error(r.Context()).Err(err).Msg("error writing text")
	}
}

func (h *httpServer) presentStatus(w http.ResponseWriter, r *http.Request, statusCode int) {
	h.logEvent(w, r, statusCode)
	w.WriteHeader(statusCode)
}

// log details of the request/response
func (h *httpServer) logEvent(w http.ResponseWriter, r *http.Request, statusCode int) {
	ctx := r.Context()
	t1 := ctx.Value(common.ContextKeyRequestStartTime).(time.Time)
	var event common.LogEvent
	if statusCode >= 500 {
		event = h.logger.Error(ctx)
	} else if statusCode >= 400 {
		event = h.logger.Warn(ctx)
	} else {
		event = h.logger.Info(ctx)
	}
	event.Strs([]struct {
		Key   string
		Value string
	}{
		{Key: "method", Value: r.Method},
		{Key: "path", Value: r.URL.Path},
		{Key: "resptime", Value: time.Since(t1).String()},
		{Key: "statuscode", Value: fmt.Sprint(statusCode)},
		{Key: "remoteaddr", Value: r.RemoteAddr},
	}).Msgf("http %d %v %v", statusCode, r.Method, r.URL.Path)
}

type bodyJson struct {
	Data     any       `json:"data"`
	NextPage *pageJson `json:"nextPage,omitempty"`
}

func toErrJson(msg string) *errJson {
	return &errJson{
		Message: msg,
	}
}

type errJson struct {
	Message string `json:"message"`
}

func (h *httpServer) getPage(r *http.Request) (pageJson, error) {
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}
	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil || offset < 0 {
		return pageJson{}, errors.New("invalid offset")
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "20"
	}
	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil || limit <= 0 {
		return pageJson{}, errors.New("invalid limit")
	}

	return pageJson{
		offset,
		limit,
		-1,
	}, nil
}

type pageJson struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
	Count  int64 `json:"count"`
}
