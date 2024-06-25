package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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
		v, err := h.metricsStorage.GetCounter(m.ID)
		if err != nil {
			http.Error(w, "counter not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		m.Value = json.RawMessage(strconv.FormatInt(v, 10))
		_ = e.Encode(m)
		return
	}

	if m.MType == "gauge" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		v, err := h.metricsStorage.GetGauge(m.ID)
		if err != nil {
			http.Error(w, "gauge not found", http.StatusNotFound)
			return
		}

		str := strings.TrimSuffix(strconv.FormatFloat(v, 'f', -1, 64), ".0")
		m.Value = json.RawMessage(str)
		_ = e.Encode(m)
		return
	}

	http.Error(w, "unknown metric type", http.StatusBadRequest)
}
