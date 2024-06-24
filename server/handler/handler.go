package handler

import "github.com/rusinov-artem/metrics/server/router"

type Handler struct {
	metricsStorage MetricsStorage
}

func New(metrics MetricsStorage) *Handler {
	h := &Handler{
		metricsStorage: metrics,
	}

	return h
}

func (h *Handler) RegisterIn(r *router.Router) *Handler {
	r.RegisterMetricsGetter(h.GetMetrics)
	r.RegisterMetricsUpdate(h.UpdateMetrics)
	r.RegisterLiveness(h.Liveness)
	r.RegisterUpdate(h.Update)
	return h
}
