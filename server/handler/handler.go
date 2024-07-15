package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/rusinov-artem/metrics/server/router"
)

type Handler struct {
	logger         *zap.Logger
	metricsStorage MetricsStorage
	dbPool         *pgxpool.Pool
}

func New(logger *zap.Logger, metrics MetricsStorage, dbpool *pgxpool.Pool) *Handler {
	h := &Handler{
		logger:         logger,
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
	r.RegisterUpdates(h.Updates)
	return h
}
