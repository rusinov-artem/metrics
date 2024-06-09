package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rusinov-artem/metrics/server/metrics"
	"github.com/rusinov-artem/metrics/server/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestUpdateMetrics(t *testing.T) {
	h := NewHandler()

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/update/counter/my_counter/42", nil)

	h.ServeHTTP(resp, req)

	b, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Result().StatusCode)
	assert.Equal(t, "", string(b))
}

type UpdateMetricsTestSuite struct {
	suite.Suite
	metrics *metrics.InMemoryMetrics
	handler http.Handler
}

func TestUpdateMetricsTestSuite(t *testing.T) {
	suite.Run(t, &UpdateMetricsTestSuite{})
}

func (s *UpdateMetricsTestSuite) SetupTest() {
	s.metrics = metrics.NewInMemory()
	handlerFn := UpdateMetrics(func() Metrics { return s.metrics })
	r := router.New()
	r.RegisterMetricsUpdate(handlerFn)
	s.handler = r.Handler()
}

func (s *UpdateMetricsTestSuite) Test_UnknwonMetricTypeError400() {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/update/unknown/my_counter/42", nil)

	s.handler.ServeHTTP(resp, req)
	s.Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *UpdateMetricsTestSuite) Test_Error404WithoutMetricName() {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/update/counter/42", nil)

	s.handler.ServeHTTP(resp, req)
	s.Equal(http.StatusNotFound, resp.Result().StatusCode)
}

func (s *UpdateMetricsTestSuite) Test_IntegerGuage() {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/integer_guage/100", nil)

	s.handler.ServeHTTP(resp, req)
	s.Equal(http.StatusOK, resp.Result().StatusCode)
}

func (s *UpdateMetricsTestSuite) Test_CanSetCounter() {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/update/counter/my_counter/42", nil)

	s.handler.ServeHTTP(resp, req)
	s.Equal(http.StatusOK, resp.Result().StatusCode)
	s.Equal(int64(42), s.metrics.Counter["my_counter"])
}

func (s *UpdateMetricsTestSuite) Test_CanSetGuage() {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/my_guage/42.42", nil)

	s.handler.ServeHTTP(resp, req)
	s.Equal(http.StatusOK, resp.Result().StatusCode)
	s.InDelta(float64(42.42), s.metrics.Guage["my_guage"], 0.0001)
}
