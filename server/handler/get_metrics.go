package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetMetrics(metricsProvider MetricsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := metricsProvider()
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")
		if metricType == "counter" {
			v, err := m.GetCounter(metricName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%d", v)
			return
		}

		if metricType == "gauge" {
			v, err := m.GetGauge(metricName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%f", v)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}
