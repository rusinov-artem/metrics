package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rusinov-artem/metrics/dto"
	serverError "github.com/rusinov-artem/metrics/server/error"
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
	var invalidRequest serverError.InvalidRequest
	if errors.As(err, &invalidRequest) {
		http.Error(w, invalidRequest.Error(), http.StatusBadRequest)
		return
	}

	var internal serverError.Internal
	if errors.As(err, &internal) {
		http.Error(w, invalidRequest.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = e.Encode(m)
}

func (h *Handler) updateSingleMetric(ctx context.Context, m *dto.Metrics) error {
	if m.MType == "counter" {
		if m.Delta == nil {
			return serverError.InvalidRequest{Msg: "counter must contain delta field"}
		}

		err := h.metricsStorage.SetCounter(ctx, m.ID, *m.Delta)
		if err != nil {
			return serverError.Internal{InnerErr: err, Msg: "unable to set counter"}
		}

		return nil
	}

	if m.MType == "gauge" {
		if m.Value == nil {
			return serverError.InvalidRequest{Msg: "gauge must contain value field"}
		}

		err := h.metricsStorage.SetGauge(ctx, m.ID, *m.Value)
		if err != nil {
			return serverError.Internal{InnerErr: err, Msg: "unable to set gauge"}
		}

		return nil
	}

	return serverError.InvalidRequest{Msg: fmt.Sprintf("unknown metric type: %s", m.MType)}
}
