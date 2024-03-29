package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/domain/usecases"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type HttpServer interface {
	Start(ctx context.Context, config HttpConfig)
	Shutdown(ctx context.Context) error
}

type httpServer struct {
	server        *http.Server
	logger        common.Logger
	config        *HttpConfig
	getChallenge  *usecases.GetChallenge
	signin        *usecases.Signin
	authenticate  *usecases.Authenticate
	rateLimit     *usecases.RateLimit
	createThread  *usecases.CreateThread
	getThread     *usecases.GetThread
	getThreads    *usecases.GetThreads
	deleteThread  *usecases.DeleteThread
	createVote    *usecases.CreateVote
	createComment *usecases.CreateComment
	getComments   *usecases.GetComments
	deleteComment *usecases.DeleteComment
	uploadImage   *usecases.UploadImage
	getUser       *usecases.GetUser
}

type HttpConfig struct {
	Port         string
	JWTSecret    string
	RealIPHeader string
}

func NewHttpServer(
	logger common.Logger,
	getChallenge *usecases.GetChallenge,
	signin *usecases.Signin,
	authenticate *usecases.Authenticate,
	rateLimit *usecases.RateLimit,
	createThread *usecases.CreateThread,
	getThread *usecases.GetThread,
	getThreads *usecases.GetThreads,
	deleteThread *usecases.DeleteThread,
	createVote *usecases.CreateVote,
	createComment *usecases.CreateComment,
	getComments *usecases.GetComments,
	deleteComment *usecases.DeleteComment,
	uploadImage *usecases.UploadImage,
	getUser *usecases.GetUser) HttpServer {
	var server *http.Server
	return &httpServer{
		server,
		logger,
		nil,
		getChallenge,
		signin,
		authenticate,
		rateLimit,
		createThread,
		getThread,
		getThreads,
		deleteThread,
		createVote,
		createComment,
		getComments,
		deleteComment,
		uploadImage,
		getUser,
	}
}

func (h *httpServer) Start(ctx context.Context, config HttpConfig) {
	h.logger.Info(ctx).Msg("starting http service")

	h.config = &config

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

	// r.Mount("/debug", middleware.Profiler())

	r.Get("/", h.healthRoute)

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.Compress(5, "application/json"))

		// public routes
		r.Group(func(r chi.Router) {
			r.Use(h.rateLimiter("public", 20, time.Minute))
			r.Use(h.maxSize(1))

			r.Get("/users/{address}", h.getUserByAddressRoute)
			r.Get("/threads", h.getThreadsRoute)
			r.Get("/threads/{threadId}", h.getThreadByIdRoute)
			r.Get("/threads/{threadId}/comments", h.getCommentsRoute)
		})

		// signin routes
		r.Group(func(r chi.Router) {
			r.Use(h.rateLimiter("signin", 5, time.Minute))
			r.Use(h.maxSize(1))

			r.Get("/signin/{address}", h.getChallengeRoute)
			r.Post("/signin/{address}", h.putChallengeRoute)
		})

		// authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(h.authentication)
			r.Use(h.ensName)
			r.Use(h.maxSize(5))

			r.With(h.rateLimiter("create:thread", 2, time.Minute*10)).Post("/threads", h.createThreadRoute)
			r.With(h.rateLimiter("vote:thread", 10, time.Minute)).Put("/threads/{threadId}/votes/{value}", h.createThreadVoteRoute)
			r.With(h.rateLimiter("create:comment", 5, time.Minute*10)).Post("/threads/{threadId}/comments", h.createCommentRoute)
			r.With(h.rateLimiter("vote:comment", 10, time.Minute)).Put("/threads/{threadId}/comments/{commentId}/votes/{value}", h.createCommentVoteRoute)
		})

		// permissioned routes
		// TODO: Make this require moderator+ permission
		r.Group(func(r chi.Router) {
			r.Use(h.authentication)
			r.Use(h.rateLimiter("permissioned", 10, time.Second))
			r.Use(h.maxSize(1))

			r.Delete("/threads/{threadId}", h.deleteThreadRoute)
			r.Delete("/threads/{threadId}/comments/{commentId}", h.deleteCommentRoute)
		})

		// image route
		r.Group(func(r chi.Router) {
			r.Use(h.authentication)
			r.Use(h.ensName)
			r.Use(h.rateLimiter("create:image", 7, time.Minute*10)) // should encompass creating images for threads and comments
			r.Use(h.maxSize(5 * 1024))

			r.Post("/images", h.uploadImageRoute)
		})
	})

	h.logger.Info(ctx).Msgf("listening on port %v", config.Port)

	h.server = &http.Server{Addr: fmt.Sprintf(":%v", config.Port), Handler: r}

	if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		h.logger.Error(ctx).Err(err).Msg("http service failed")
		panic(err)
	}

	h.logger.Info(ctx).Msg("http service stopped")
}

func (h *httpServer) Shutdown(ctx context.Context) error {
	h.logger.Info(ctx).Msg("shutting down http service")

	if h.server != nil {
		return h.server.Shutdown(ctx)
	}
	return nil
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

func (h *httpServer) presentForbidden(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Warn(r.Context()).Err(err).Msg("forbidden")
	h.presentJSON(w, r, http.StatusForbidden, toErrJson("forbidden"), nil)
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
	routeCtx := chi.RouteContext(ctx)
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
		{Key: "pattern", Value: fmt.Sprintf("%v %v", routeCtx.RouteMethod, routeCtx.RoutePattern())},
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
