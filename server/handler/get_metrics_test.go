package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/metrics/server/metrics"
	"github.com/rusinov-artem/metrics/server/router"
)

type MetricsGetterTestSuite struct {
	suite.Suite
	metrics *metrics.InMemoryMetrics
	handler http.Handler
}

func TestMetricsGetter(t *testing.T) {
	suite.Run(t, &MetricsGetterTestSuite{})
}

func (s *MetricsGetterTestSuite) SetupTest() {
	s.metrics = metrics.NewInMemory()
	handlerFn := New(s.metrics).GetMetrics
	r := router.New()
	r.RegisterMetricsGetter(handlerFn)
	s.handler = r.Mux()
}

func (s *MetricsGetterTestSuite) Test_ErrorIfCounterNameNotFound() {
	req := httptest.NewRequest(http.MethodGet, "/value/counter/unknown", nil)
	res := s.Do(req)
	defer closeBody(res)

	s.Equal(http.StatusNotFound, res.StatusCode)
}

func (s *MetricsGetterTestSuite) Test_ErrorIfGaugeNameNotFound() {
	req := httptest.NewRequest(http.MethodGet, "/value/gauge/unknown", nil)
	res := s.Do(req)
	defer closeBody(res)

	s.Equal(http.StatusNotFound, res.StatusCode)
}

func (s *MetricsGetterTestSuite) Test_CanGetCounterValue() {
	s.metrics.Counter["my_counter"] = 42
	req := httptest.NewRequest(http.MethodGet, "/value/counter/my_counter", nil)
	resp := s.Do(req)
	defer closeBody(resp)
	s.Equal(http.StatusOK, resp.StatusCode)
	b, _ := io.ReadAll(resp.Body)
	v, err := strconv.ParseInt(string(b), 10, 64)
	s.NoError(err)
	s.Equal(int64(42), v)
}

func (s *MetricsGetterTestSuite) Test_CanGetGaugeValue() {
	s.metrics.Gauge["my_gauge"] = 42.42
	req := httptest.NewRequest(http.MethodGet, "/value/gauge/my_gauge", nil)
	resp := s.Do(req)
	defer closeBody(resp)
	s.Equal(http.StatusOK, resp.StatusCode)
	b, _ := io.ReadAll(resp.Body)
	v, err := strconv.ParseFloat(string(b), 64)
	s.NoError(err)
	s.InDelta(v, 42.42, 0.0001)
}

func (s *MetricsGetterTestSuite) Test_BadRequest() {
	req := httptest.NewRequest(http.MethodGet, "/value/unknown/my_gauge", nil)
	resp := s.Do(req)
	defer closeBody(resp)
	s.Equal(http.StatusBadRequest, resp.StatusCode)
	b, err := io.ReadAll(resp.Body)
	s.NoError(err)
	s.Empty(b)
}

func (s *MetricsGetterTestSuite) Do(req *http.Request) *http.Response {
	recorder := httptest.NewRecorder()
	s.handler.ServeHTTP(recorder, req)
	return recorder.Result()
}
