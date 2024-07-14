package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MetricsStorage interface {
	SetCounter(ctx context.Context, name string, value int64) error
	SetGauge(ctx context.Context, name string, value float64) error

	GetCounter(ctx context.Context, name string) (int64, error)
	GetGauge(ctx context.Context, name string) (float64, error)
}

func (h *Handler) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	if metricType == "counter" {
		v, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = h.metricsStorage.SetCounter(r.Context(), metricName, v)
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
		_ = h.metricsStorage.SetGauge(r.Context(), metricName, v)
		w.WriteHeader(http.StatusOK)
		log.Printf("updated gauge '%s' = %f", metricName, v)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
