package router

import (
	"net/http"

	"github.com/rusinov-artem/metrics/server/middleware"
)

type Mux interface {
	Method(method, pattern string, handler http.Handler)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Handler interface {
	GetMetrics(http.ResponseWriter, *http.Request)
	UpdateMetrics(http.ResponseWriter, *http.Request)
	Liveness(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	Value(http.ResponseWriter, *http.Request)
	Info(http.ResponseWriter, *http.Request)
	Ping(http.ResponseWriter, *http.Request)
	Updates(http.ResponseWriter, *http.Request)
}

type Router struct {
	mux        Mux
	middleware middleware.Middleware
}

func New(m Mux) *Router {
	return &Router{
		mux:        m,
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

func (t *Router) RegisterUpdates(h http.HandlerFunc) {
	t.mux.Method(http.MethodPost, "/updates/", h)
}

func (t *Router) withMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.middleware(h).ServeHTTP(w, r)
	})
}

func (t *Router) Mux() http.Handler {
	return t.withMiddleware(t.mux)
}

func (t *Router) SetHandler(h Handler) *Router {
	t.RegisterMetricsGetter(h.GetMetrics)
	t.RegisterMetricsUpdate(h.UpdateMetrics)
	t.RegisterLiveness(h.Liveness)
	t.RegisterUpdate(h.Update)
	t.RegisterValue(h.Value)
	t.RegisterInfo(h.Info)
	t.RegisterPing(h.Ping)
	t.RegisterUpdates(h.Updates)
	return t
}
