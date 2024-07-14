package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	if metricType == "counter" {
		v, err := h.metricsStorage.GetCounter(r.Context(), metricName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "%d", v)
		return
	}

	if metricType == "gauge" {
		v, err := h.metricsStorage.GetGauge(r.Context(), metricName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		stringV := strings.TrimRight(fmt.Sprintf("%.3f", v), "0.")
		_, _ = fmt.Fprint(w, stringV)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
