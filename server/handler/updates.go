package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rusinov-artem/metrics/dto"
)

func (h *Handler) Updates(w http.ResponseWriter, r *http.Request) {
	m := &[]dto.Metrics{}
	d := json.NewDecoder(r.Body)
	if err := d.Decode(m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range *m {
		err := h.updateSingleMetric(r.Context(), &(*m)[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
