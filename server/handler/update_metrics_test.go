package handler

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/metrics/server/metrics"
	"github.com/rusinov-artem/metrics/server/router"
)

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
	handlerFn := New(s.metrics).UpdateMetrics
	r := router.New()
	r.RegisterMetricsUpdate(handlerFn)
	s.handler = r.Mux()
}

func (s *UpdateMetricsTestSuite) Do(req *http.Request) *http.Response {
	recorder := httptest.NewRecorder()
	s.handler.ServeHTTP(recorder, req)
	return recorder.Result()
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
	s.InDelta(float64(42.42), s.metrics.Gauge["my_guage"], 0.0001)
}

func (s *UpdateMetricsTestSuite) Test_Race() {
	req1 := httptest.NewRequest(http.MethodPost, "/update/gauge/my_guage1/42.42", nil)
	req2 := httptest.NewRequest(http.MethodPost, "/update/gauge/my_guage2/47.47", nil)

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
		defer func() {
			_ = resp.Body.Close()
		}()
	}
}
