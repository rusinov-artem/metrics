package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Metrics interface {
	SetCounter(name string, value int64) error
	SetGuage(name string, value float64) error
}

type MetricsProvider func() Metrics

func UpdateMetrics(metricsProvider MetricsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := metricsProvider()
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")
		if metricType == "counter" {
			v, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
			if err != nil {
				log.Println(err)

				w.WriteHeader(http.StatusBadRequest)
				return
			}
			_ = m.SetCounter(metricName, v)
			w.WriteHeader(http.StatusOK)
			return
		}

		if metricType == "gauge" {
			v, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
			if err != nil {
				log.Println(err)

				w.WriteHeader(http.StatusBadRequest)
				return
			}
			_ = m.SetGuage(metricName, v)
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}
