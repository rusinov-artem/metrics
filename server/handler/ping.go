package handler

import (
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancelFN := h.context(r.Context())
	defer cancelFN()

	if h.dbPool == nil {
		h.logger.Error("db pool not initialized")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := h.dbPool.Ping(ctx)
	if err != nil {
		h.logger.Error("ping error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
