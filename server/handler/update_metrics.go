package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MetricsStorage interface {
	SetCounter(name string, value int64) error
	SetGauge(name string, value float64) error
	Flush(ctx context.Context) error

	GetCounter(ctx context.Context, name string) (int64, error)
	GetGauge(ctx context.Context, name string) (float64, error)
}

func (h *Handler) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	ctx, cancelFN := h.context(r.Context())
	defer cancelFN()

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	if metricType == "counter" {
		v, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = h.metricsStorage.SetCounter(metricName, v)
		_ = h.metricsStorage.Flush(ctx)
		w.WriteHeader(http.StatusOK)
		log.Printf("updated counter '%s' = %d", metricName, v)
		return
	}

	if metricType == "gauge" {
		v, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = h.metricsStorage.SetGauge(metricName, v)
		_ = h.metricsStorage.Flush(ctx)
		w.WriteHeader(http.StatusOK)
		log.Printf("updated gauge '%s' = %f", metricName, v)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
