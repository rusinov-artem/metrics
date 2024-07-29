package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rusinov-artem/metrics/dto"
	serverError "github.com/rusinov-artem/metrics/server/error"
)

func (h *Handler) Updates(w http.ResponseWriter, r *http.Request) {
	ctx, cancelFN := h.context(r.Context())
	defer cancelFN()

	m := &[]dto.Metrics{}
	d := json.NewDecoder(r.Body)
	if err := d.Decode(m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range *m {
		err := h.updateSingleMetric(&(*m)[i])
		var internalError serverError.Internal
		if errors.As(err, &internalError) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var invalidRequest serverError.InvalidRequest
		if errors.As(err, &invalidRequest) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err := h.metricsStorage.Flush(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
