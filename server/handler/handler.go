package handler

import "github.com/rusinov-artem/metrics/server/router"

type Handler struct {
	metrics Metrics
}

func New(metrics Metrics) *Handler {
	h := &Handler{
		metrics: metrics,
	}

	return h
}

func (h *Handler) RegisterIn(r *router.Router) *Handler {
	r.RegisterMetricsGetter(h.GetMetrics)
	r.RegisterMetricsUpdate(h.UpdateMetrics)
	r.RegisterLiveness(h.Liveness)
	return h
}
