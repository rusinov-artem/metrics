package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rusinov-artem/metrics/server/metrics"
	"github.com/rusinov-artem/metrics/server/router"
	"github.com/stretchr/testify/suite"
)

type UpdateMetricsTestSuite struct {
	suite.Suite
	metrics  *metrics.InMemoryMetrics
	handler  http.Handler
	recorder *httptest.ResponseRecorder
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
	s.recorder = httptest.NewRecorder()
}

func (s *UpdateMetricsTestSuite) Do(req *http.Request) *http.Response {
	s.handler.ServeHTTP(s.recorder, req)
	return s.recorder.Result()
}

func (s *UpdateMetricsTestSuite) Test_UnknwonMetricTypeError400() {
	req := httptest.NewRequest(http.MethodPost, "/update/unknown/my_counter/42", nil)

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *UpdateMetricsTestSuite) Test_Error404WithoutMetricName() {
	req := httptest.NewRequest(http.MethodPost, "/update/counter/42", nil)

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *UpdateMetricsTestSuite) Test_IntegerGuage() {
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/integer_guage/100", nil)

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusOK, resp.StatusCode)
}

func (s *UpdateMetricsTestSuite) Test_CanSetCounter() {
	req := httptest.NewRequest(http.MethodPost, "/update/counter/my_counter/42", nil)

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(int64(42), s.metrics.Counter["my_counter"])
}

func (s *UpdateMetricsTestSuite) Test_CanSetGuage() {
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/my_guage/42.42", nil)

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.InDelta(float64(42.42), s.metrics.Guage["my_guage"], 0.0001)
}
