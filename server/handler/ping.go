package handler

import (
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if h.dbPool == nil {
		h.logger.Error("db pool not initialized")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := h.dbPool.Ping(r.Context())
	if err != nil {
		h.logger.Error("ping error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
