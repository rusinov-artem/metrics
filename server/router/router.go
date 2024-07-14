package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/rusinov-artem/metrics/server/middleware"
)

type Router struct {
	mux        *chi.Mux
	middleware middleware.Middleware
}

func New() *Router {
	mux := chi.NewRouter()
	return &Router{
		mux:        mux,
		middleware: middleware.Nop,
	}
}

func (t *Router) AddMiddleware(m middleware.Middleware) {
	t.middleware = t.middleware.Wrap(m)
}

func (t *Router) RegisterLiveness(h http.HandlerFunc) {
	t.mux.Method(http.MethodGet, "/liveness", t.withMiddleware(h))
}

func (t *Router) RegisterMetricsUpdate(h http.HandlerFunc) {
	t.mux.Method(http.MethodPost, "/update/{type}/{name}/{value}", h)
}

func (t *Router) RegisterMetricsGetter(h http.HandlerFunc) {
	t.mux.Method(http.MethodGet, "/value/{type}/{name}", h)
}

func (t *Router) RegisterUpdate(h http.HandlerFunc) {
	t.mux.Method(http.MethodPost, "/update/", h)
}

func (t *Router) RegisterValue(h http.HandlerFunc) {
	t.mux.Method(http.MethodPost, "/value/", h)
}

func (t *Router) RegisterInfo(h http.HandlerFunc) {
	t.mux.Method(http.MethodGet, "/", h)
}

func (t *Router) RegisterPing(h http.HandlerFunc) {
	t.mux.Method(http.MethodGet, "/ping", h)
}

func (t *Router) withMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.middleware(h).ServeHTTP(w, r)
	})
}

func (t *Router) Mux() http.Handler {
	return t.withMiddleware(t.mux)
}
