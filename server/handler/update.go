package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/rusinov-artem/metrics/dto"
	serverError "github.com/rusinov-artem/metrics/server/error"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancelFN := h.context(r.Context())
	defer cancelFN()

	m := &dto.Metrics{}
	d := json.NewDecoder(r.Body)
	e := json.NewEncoder(w)
	if err := d.Decode(m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.updateSingleMetric(m)
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

	err = h.metricsStorage.Flush(ctx)
	if err != nil {
		h.logger.Error("unable to flush storage", zap.Error(err))
		http.Error(w, invalidRequest.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = e.Encode(m)
}

func (h *Handler) updateSingleMetric(m *dto.Metrics) error {
	if m.MType == "counter" {
		if m.Delta == nil {
			return serverError.InvalidRequest{Msg: "counter must contain delta field"}
		}

		err := h.metricsStorage.SetCounter(m.ID, *m.Delta)
		if err != nil {
			return serverError.Internal{InnerErr: err, Msg: "unable to set counter"}
		}

		return nil
	}

	if m.MType == "gauge" {
		if m.Value == nil {
			return serverError.InvalidRequest{Msg: "gauge must contain value field"}
		}

		err := h.metricsStorage.SetGauge(m.ID, *m.Value)
		if err != nil {
			return serverError.Internal{InnerErr: err, Msg: "unable to set gauge"}
		}

		return nil
	}

	return serverError.InvalidRequest{Msg: fmt.Sprintf("unknown metric type: %s", m.MType)}
}
