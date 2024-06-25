package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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
		if len(m.Value) < 1 {
			http.Error(w, "counter value must be set", http.StatusBadRequest)
			return
		}

		v, err := strconv.ParseInt(string(m.Value), 10, 64)
		if err != nil {
			http.Error(w, "unable to parse value", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = h.metricsStorage.SetCounter(m.ID, v)
		_ = e.Encode(m)
		return
	}

	if m.MType == "gauge" {
		if len(m.Value) < 1 {
			http.Error(w, "counter value must be set", http.StatusBadRequest)
			return

		}

		v, err := strconv.ParseFloat(string(m.Value), 64)
		if err != nil {
			http.Error(w, "unable to parse value", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = h.metricsStorage.SetGauge(m.ID, v)
		_ = e.Encode(m)
		return
	}

	http.Error(w, "unknown metric type", http.StatusBadRequest)
}
