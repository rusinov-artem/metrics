package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/rusinov-artem/metrics/dto"
	"github.com/rusinov-artem/metrics/server/middleware"
	"github.com/rusinov-artem/metrics/server/router"
	"github.com/rusinov-artem/metrics/server/storage"
)

type UpdateTestSuite struct {
	suite.Suite
	metrics *storage.InMemoryMetrics
	handler http.Handler
}

func TestUpdateTestSuite(t *testing.T) {
	suite.Run(t, &UpdateTestSuite{})
}

func (s *UpdateTestSuite) SetupTest() {
	s.metrics = storage.NewInMemory()
	metricsFactory := func() MetricsStorage {
		return s.metrics
	}
	handlerFn := New(nil, metricsFactory, nil).Update
	r := router.New()
	r.AddMiddleware(middleware.Logger(zap.NewNop()))
	r.AddMiddleware(middleware.GzipEncoder())
	r.RegisterUpdate(handlerFn)
	s.handler = r.Mux()
}

func (s *UpdateTestSuite) Test_ErrorOnBadRequest() {
	req := httptest.NewRequest(http.MethodPost, "/update/", nil)

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *UpdateTestSuite) Test_Error400WithoutType() {
	req := httptest.NewRequest(http.MethodPost, "/update/", Empty())

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *UpdateTestSuite) Test_Error400WithoutCounterValue() {
	req := httptest.NewRequest(http.MethodPost, "/update/", CounterNoValue())

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *UpdateTestSuite) Test_CanUpdateCounterWithValue() {
	req := httptest.NewRequest(http.MethodPost, "/update/", counter(42))

	resp := s.Do(req)
	defer closeBody(resp)

	actual := dto.Metrics{}
	d := json.NewDecoder(resp.Body)
	err := d.Decode(&actual)
	s.Require().NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(int64(42), *actual.Delta)
	s.Equal("application/json", resp.Header.Get("Content-Type"))
}

func (s *UpdateTestSuite) Test_CanUpdateCounterWithNameAndValue() {
	req := httptest.NewRequest(http.MethodPost, "/update/", counterWithName("my_counter", 42))

	resp := s.Do(req)
	defer closeBody(resp)

	actual := dto.Metrics{}
	d := json.NewDecoder(resp.Body)
	err := d.Decode(&actual)
	s.Require().NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal("my_counter", actual.ID)
	s.Equal(int64(42), *actual.Delta)
	s.Equal("application/json", resp.Header.Get("Content-Type"))
}

func (s *UpdateTestSuite) Test_Error400WithoutGaugeValue() {
	req := httptest.NewRequest(http.MethodPost, "/update/", GaugeNoValue())

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *UpdateTestSuite) Test_CanUpdateGaugeWithValue() {
	req := httptest.NewRequest(http.MethodPost, "/update/", gauge(42.42))

	resp := s.Do(req)
	defer closeBody(resp)

	actual := dto.Metrics{}
	d := json.NewDecoder(resp.Body)
	err := d.Decode(&actual)
	s.Require().NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.InDelta(42.42, *actual.Value, 0.001)
	s.Equal("application/json", resp.Header.Get("Content-Type"))
}

func (s *UpdateTestSuite) Test_CanUpdateGaugeWithNameAndValue() {
	req := httptest.NewRequest(http.MethodPost, "/update/", gaugeWithName("my_gauge", 42.42))

	resp := s.Do(req)
	defer closeBody(resp)

	actual := dto.Metrics{}
	d := json.NewDecoder(resp.Body)
	err := d.Decode(&actual)
	s.Require().NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.InDelta(42.42, *actual.Value, 0.001)
	s.Equal("my_gauge", actual.ID)
	s.Equal("application/json", resp.Header.Get("Content-Type"))
}

func (s *UpdateTestSuite) Test_Race() {
	req1 := httptest.NewRequest(http.MethodPost, "/update/", gaugeWithName("g1", 13.13))
	req2 := httptest.NewRequest(http.MethodPost, "/update/", gaugeWithName("g2", 14.13))

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

func (s *UpdateTestSuite) ExecuteTimes(times int, req *http.Request) {
	for i := 0; i < times; i++ {
		resp := s.Do(req)
		_ = resp.Body.Close()
	}
}

func (s *UpdateTestSuite) Do(req *http.Request) *http.Response {
	recorder := httptest.NewRecorder()
	s.handler.ServeHTTP(recorder, req)
	return recorder.Result()
}

func counter(v int) io.Reader {
	return bytes.NewBufferString(
		fmt.Sprintf(`
	{
		"type": "counter",
		"delta": %d
	}
`, v))
}

func gauge(v float64) io.Reader {
	return bytes.NewBufferString(
		fmt.Sprintf(`
	{
		"type": "gauge",
		"value": %f
	}
`, v))
}

func Empty() io.Reader {
	return bytes.NewBufferString(`
	{}
`)
}

func CounterNoValue() io.Reader {
	return bytes.NewBufferString(`
	{
		"type": "counter"
	}
`)
}

func GaugeNoValue() io.Reader {
	return bytes.NewBufferString(`
	{
		"type": "gauge"
	}
`)
}

func counterWithName(n string, v int) io.Reader {
	return bytes.NewBufferString(
		fmt.Sprintf(`
	{
		"id":"%s",
		"type": "counter",
		"delta": %d
	}
`, n, v))
}

func gaugeWithName(n string, v float64) io.Reader {
	return bytes.NewBufferString(
		fmt.Sprintf(`
	{
		"id":"%s",
		"type": "gauge",
		"value": %f
	}
`, n, v))
}
