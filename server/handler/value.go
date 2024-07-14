package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rusinov-artem/metrics/dto"
)

func (h *Handler) Value(w http.ResponseWriter, r *http.Request) {
	m := &dto.Metrics{}
	d := json.NewDecoder(r.Body)
	e := json.NewEncoder(w)
	if err := d.Decode(m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if m.MType == "counter" {
		v, err := h.metricsStorage.GetCounter(r.Context(), m.ID)
		if err != nil {
			http.Error(w, "counter not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		m.Delta = &v
		_ = e.Encode(m)
		return
	}

	if m.MType == "gauge" {

		v, err := h.metricsStorage.GetGauge(r.Context(), m.ID)
		if err != nil {
			http.Error(w, "gauge not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		m.Value = &v
		_ = e.Encode(m)
		return
	}

	http.Error(w, "unknown metric type", http.StatusBadRequest)
}
