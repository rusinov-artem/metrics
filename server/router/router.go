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

func (this *Router) RegisterLiveness(h http.HandlerFunc) {
	this.mux.Get("/liveness", h)
}

func (this *Router) RegisterMetricsUpdate(h http.HandlerFunc) {
	this.mux.Post("/update/{type}/{name}/{value}", h)
}

func (this *Router) Handler() http.Handler {
	return this.mux
}
