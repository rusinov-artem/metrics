package handler

import (
	"context"
	"encoding/json"
	"fmt"
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

	err := h.updateSingleMetric(r.Context(), m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = e.Encode(m)
}

func (h *Handler) updateSingleMetric(ctx context.Context, m *dto.Metrics) error {
	if m.MType == "counter" {
		if m.Delta == nil {
			return fmt.Errorf("counter value must be set")
		}

		_ = h.metricsStorage.SetCounter(ctx, m.ID, *m.Delta)
		return nil
	}

	if m.MType == "gauge" {
		if m.Value == nil {
			return fmt.Errorf("counter value must be set")
		}

		_ = h.metricsStorage.SetGauge(ctx, m.ID, *m.Value)
		return nil
	}

	return fmt.Errorf("unknown metric type: %s", m.MType)
}
