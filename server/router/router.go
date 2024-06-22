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
	t.mux.Method(http.MethodGet, "/liveness", t.middleware(h))
}

func (t *Router) RegisterMetricsUpdate(h http.HandlerFunc) {
	t.mux.Method(http.MethodPost, "/update/{type}/{name}/{value}", t.middleware(h))
}

func (t *Router) RegisterMetricsGetter(h http.HandlerFunc) {
	t.mux.Method(http.MethodGet, "/value/{type}/{name}", t.middleware(h))
}

func (t *Router) Mux() http.Handler {
	return t.mux
}
