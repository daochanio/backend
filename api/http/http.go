package http

import (
	"context"
	"encoding/json"
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

type IHttpServer interface {
	Start(context.Context) error
}

type httpServer struct {
	logger              common.ILogger
	settings            settings.ISettings
	createThreadUseCase *usecases.CreateThreadUseCase
	getThreadUseCase    *usecases.GetThreadUseCase
	getThreadsUseCase   *usecases.GetThreadsUseCase
	deleteThreadUseCase *usecases.DeleteThreadUseCase
	voteThreadUseCase   *usecases.VoteThreadUseCase
}

func NewHttpServer(
	logger common.ILogger,
	settings settings.ISettings,
	createThreadUseCase *usecases.CreateThreadUseCase,
	getThreadUseCase *usecases.GetThreadUseCase,
	getThreadsUseCase *usecases.GetThreadsUseCase,
	deleteThreadUseCase *usecases.DeleteThreadUseCase,
	voteThreadUseCase *usecases.VoteThreadUseCase) IHttpServer {
	return &httpServer{
		logger,
		settings,
		createThreadUseCase,
		getThreadUseCase,
		getThreadsUseCase,
		deleteThreadUseCase,
		voteThreadUseCase,
	}
}

func (h *httpServer) Start(ctx context.Context) error {
	h.logger.Info(ctx).Msg("starting http service")

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://daochan.io", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.NoCache)
	r.Use(h.timer)
	r.Use(h.realIP)
	r.Use(h.requestId)
	r.Use(h.recoverer)
	r.Use(h.timeout)

	r.Get("/", h.healthRoute)

	r.Route("/threads", func(r chi.Router) {
		r.Post("/", h.createThreadRoute)
		r.Get("/", h.getThreadsRoute)
		r.Get("/{id}", h.getThreadByIdRoute)
		r.Delete("/{id}", h.deleteThreadRoute)
		r.Put("/{id}/vote/{vote}", h.voteThreadRoute)
	})

	port := h.settings.Port()

	h.logger.Info(ctx).Msgf("listening on port %v", port)

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), r)

	h.logger.Error(ctx).Err(err).Msg("error in http service")

	return err
}

func (h *httpServer) presentNotFound(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Info(r.Context()).Err(err).Msg("not found")
	h.presentJSON(w, r, http.StatusNotFound, toErrJson("not found"))
}

func (h *httpServer) presentBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Info(r.Context()).Err(err).Msg("bad request")
	h.presentJSON(w, r, http.StatusBadRequest, toErrJson("bad request"))
}

func (h *httpServer) presentJSON(w http.ResponseWriter, r *http.Request, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	h.presentStatus(w, r, statusCode)
	json.NewEncoder(w).Encode(body)
}

func (h *httpServer) presentText(w http.ResponseWriter, r *http.Request, statusCode int, text string) {
	w.Header().Set("Content-Type", "text/plain")
	h.presentStatus(w, r, statusCode)
	w.Write([]byte(text))
}

func (h *httpServer) presentStatus(w http.ResponseWriter, r *http.Request, statusCode int) {
	h.logEvent(w, r, statusCode)
	w.WriteHeader(statusCode)
}

// log details of the request/response
func (h *httpServer) logEvent(w http.ResponseWriter, r *http.Request, statusCode int) {
	ctx := r.Context()
	t1 := ctx.Value(common.ContextKeyRequestStartTime).(time.Time)
	h.logger.Info(ctx).Strs([]struct {
		Key   string
		Value string
	}{
		{Key: "method", Value: r.Method},
		{Key: "path", Value: r.URL.Path},
		{Key: "resptime", Value: time.Since(t1).String()},
		{Key: "statuscode", Value: fmt.Sprint(statusCode)},
		{Key: "remoteaddr", Value: r.RemoteAddr},
	}).Msg("http event")
}

func toErrJson(msg string) *errJson {
	return &errJson{
		Message: msg,
	}
}

type errJson struct {
	Message string `json:"message"`
}

func (h *httpServer) getPaginationParams(r *http.Request) (paginationParams, error) {
	offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 32)
	if err != nil {
		offset = 0
	}

	limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 32)
	if err != nil {
		limit = 100
	}

	return paginationParams{
		uint32(offset),
		uint32(limit),
	}, nil
}

type paginationParams struct {
	Offset uint32
	Limit  uint32
}
