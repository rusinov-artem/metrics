package handler

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/rusinov-artem/metrics/server/middleware"
	"github.com/rusinov-artem/metrics/server/router"
	"github.com/rusinov-artem/metrics/server/storage"
)

type UpdateMetricsTestSuite struct {
	suite.Suite
	metrics *storage.InMemoryMetrics
	handler http.Handler
}

func TestUpdateMetricsTestSuite(t *testing.T) {
	suite.Run(t, &UpdateMetricsTestSuite{})
}

func (s *UpdateMetricsTestSuite) SetupTest() {
	s.metrics = storage.NewInMemory()
	handlerFn := New(nil, s.metrics, nil).UpdateMetrics
	r := router.New()
	r.AddMiddleware(middleware.Logger(zap.NewNop()))
	r.RegisterMetricsUpdate(handlerFn)
	s.handler = r.Mux()
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

func (s *UpdateMetricsTestSuite) Test_ErrorCounterWithBadValue() {
	req := httptest.NewRequest(http.MethodPost, "/update/counter/name/bad_value", nil)

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *UpdateMetricsTestSuite) Test_ErrorGaugeWithBadValue() {
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/name/bad_value", nil)

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *UpdateMetricsTestSuite) Test_IntegerGauge() {
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/integer_gauge/100", nil)

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

func (s *UpdateMetricsTestSuite) Test_CanSetGauge() {
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/my_gauge/42.42", nil)

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.InDelta(42.42, s.metrics.Gauge["my_gauge"], 0.0001)
}

func (s *UpdateMetricsTestSuite) Test_Race() {
	req1 := httptest.NewRequest(http.MethodPost, "/update/gauge/my_gauge1/42.42", nil)
	req2 := httptest.NewRequest(http.MethodPost, "/update/gauge/my_gauge2/47.47", nil)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		s.ExecuteTimes(100, req1)
	}()

	go func() {
		defer wg.Done()
		s.ExecuteTimes(100, req2)
	}()

	wg.Wait()

}

func (s *UpdateMetricsTestSuite) ExecuteTimes(times int, req *http.Request) {
	for i := 0; i < times; i++ {
		resp := s.Do(req)
		_ = resp.Body.Close()
	}
}

func (s *UpdateMetricsTestSuite) Do(req *http.Request) *http.Response {
	recorder := httptest.NewRecorder()
	s.handler.ServeHTTP(recorder, req)
	return recorder.Result()
}
