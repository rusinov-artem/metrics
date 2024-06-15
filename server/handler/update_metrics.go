package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Metrics interface {
	SetCounter(name string, value int64) error
	SetGauge(name string, value float64) error

	GetCounter(name string) (int64, error)
	GetGauge(name string) (float64, error)
}

type MetricsProvider func() Metrics

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
		_ = h.metrics.SetCounter(metricName, v)
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
		_ = h.metrics.SetGauge(metricName, v)
		w.WriteHeader(http.StatusOK)
		log.Printf("updated gauge '%s' = %f", metricName, v)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
