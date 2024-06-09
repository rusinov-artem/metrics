package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	mux *chi.Mux
}

func New() *Router {
	mux := chi.NewRouter()
	return &Router{
		mux: mux,
	}
}

func (t *Router) RegisterLiveness(h http.HandlerFunc) {
	t.mux.Get("/liveness", h)
}

func (t *Router) RegisterMetricsUpdate(h http.HandlerFunc) {
	t.mux.Post("/update/{type}/{name}/{value}", h)
}

func (t *Router) Handler() http.Handler {
	return t.mux
}
