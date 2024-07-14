package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rusinov-artem/metrics/dto"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	m := &dto.Metrics{}
	d := json.NewDecoder(r.Body)
	e := json.NewEncoder(w)
	if err := d.Decode(m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if m.MType == "counter" {
		if m.Delta == nil {
			http.Error(w, "counter value must be set", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = h.metricsStorage.SetCounter(r.Context(), m.ID, *m.Delta)
		_ = e.Encode(m)
		return
	}

	if m.MType == "gauge" {
		if m.Value == nil {
			http.Error(w, "counter value must be set", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = h.metricsStorage.SetGauge(r.Context(), m.ID, *m.Value)
		_ = e.Encode(m)
		return
	}

	http.Error(w, "unknown metric type", http.StatusBadRequest)
}
