package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rusinov-artem/metrics/server/router"
)

type Handler struct {
	metricsStorage MetricsStorage
	dbPool         *pgxpool.Pool
}

func New(metrics MetricsStorage, dbpool *pgxpool.Pool) *Handler {
	h := &Handler{
		metricsStorage: metrics,
		dbPool:         dbpool,
	}

	return h
}

func (h *Handler) RegisterIn(r *router.Router) *Handler {
	r.RegisterMetricsGetter(h.GetMetrics)
	r.RegisterMetricsUpdate(h.UpdateMetrics)
	r.RegisterLiveness(h.Liveness)
	r.RegisterUpdate(h.Update)
	r.RegisterValue(h.Value)
	r.RegisterInfo(h.Info)
	r.RegisterPing(h.Ping)
	return h
}
