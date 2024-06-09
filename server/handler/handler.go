package handler

import (
	"net/http"

	"github.com/rusinov-artem/metrics/server/metrics"
	"github.com/rusinov-artem/metrics/server/router"
)

func NewHandler() http.Handler {

	m := metrics.NewInMemory()
	metricsProvider := func() Metrics {
		return m
	}

	r := router.New()
	r.RegisterLiveness(Liveness())
	r.RegisterMetricsUpdate(UpdateMetrics(metricsProvider))
	r.RegisterMetricsGetter(GetMetrics(metricsProvider))

	return r.Handler()
}
