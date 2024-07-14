package handler

import "net/http"

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.dbPool.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
