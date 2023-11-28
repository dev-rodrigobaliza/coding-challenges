package api

import (
	"context"
	"log/slog"
	"net/http"
	"rl/limiter"
	"rl/log"
	"time"

	"github.com/google/uuid"
)

type ContextKey string

const ContextKeyID ContextKey = "id"

type Api struct {
	rl     *limiter.Limiter
	logger *slog.Logger
	server *http.Server
}

func New(algo limiter.Algorithm, logAll bool) *Api {
	rl := limiter.New(algo)
	logger := log.New(logAll)

	mux := http.NewServeMux()

	server := http.Server{
		Handler: mux,
	}

	api := Api{
		rl:     rl,
		logger: logger,
		server: &server,
	}

	mux.Handle("/limited", api.loggingMiddleware(api.ratelimitMiddleware(http.HandlerFunc(getLimited))))
	mux.Handle("/unlimited", api.loggingMiddleware(http.HandlerFunc(getUnlimited)))

	return &api
}

func (a *Api) Start(addr string) {
	a.server.Addr = addr
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	a.logger.Info("server started", slog.String("addr", addr))
}

func (a *Api) Stop() {
	if err := a.server.Shutdown(context.Background()); err != nil {
		a.logger.Error("failed to close api server", slog.Any("err", err))
	}

	a.rl.Stop()
	a.logger.Info("server stopped")
}

func (a *Api) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		id := uuid.New().String()
		a.logger.Info("request started", slog.String("id", id), slog.String("method", r.Method), slog.String("path", r.URL.Path), slog.String("client", r.RemoteAddr))
		ctx := context.WithValue(r.Context(), ContextKeyID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
		a.logger.Info("request finished", slog.String("id", id), slog.String("elapsed", time.Since(start).String()))
	})
}

func (a *Api) ratelimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(ContextKeyID)
		ip := r.RemoteAddr
		if a.rl.Can(id.(string), ip) {
			next.ServeHTTP(w, r)
			return
		}

		a.logger.Info("rate limited request", slog.String("method", r.Method), slog.String("path", r.URL.Path), slog.String("client", r.RemoteAddr))
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
	})
}

func getLimited(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Limited, don't over use me!"))
}

func getUnlimited(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Unlimited! Let's Go!"))
}
