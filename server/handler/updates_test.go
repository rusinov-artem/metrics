package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/rusinov-artem/metrics/server/middleware"
	"github.com/rusinov-artem/metrics/server/router"
	"github.com/rusinov-artem/metrics/server/storage"
)

type UpdatesTestSuite struct {
	suite.Suite
	metrics *storage.InMemoryMetrics
	handler http.Handler
}

func TestUpdatesTestSuite(t *testing.T) {
	suite.Run(t, &UpdatesTestSuite{})
}

func (s *UpdatesTestSuite) SetupTest() {
	s.metrics = storage.NewInMemory()
	handlerFn := New(nil, s.metrics, nil).Updates
	r := router.New()
	r.AddMiddleware(middleware.Logger(zap.NewNop()))
	r.AddMiddleware(middleware.GzipEncoder())
	r.RegisterUpdate(handlerFn)
	s.handler = r.Mux()
}

func (s *UpdatesTestSuite) Test_ErrorOnInvalidJson() {
	req := httptest.NewRequest(http.MethodPost, "/update/", nil)

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *UpdatesTestSuite) Test_CanHandleBatchOfSingleElement() {
	req := httptest.NewRequest(http.MethodPost, "/update/", singleElementBatch())

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(int64(3), s.metrics.Counter["single_element"])
}

func (s *UpdatesTestSuite) Test_CanHandleMultiElementBatch() {
	req := httptest.NewRequest(http.MethodPost, "/update/", multiElementBatch())

	resp := s.Do(req)
	defer closeBody(resp)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(int64(42), s.metrics.Counter["first_counter"])
	s.InDelta(43.44, s.metrics.Gauge["first_gauge"], 0.0001)
}

func multiElementBatch() io.Reader {
	return bytes.NewBufferString(`
	[
		{
			"id":"first_counter",
			"type": "counter",
			"delta": 42
	    },
		{
			"id":"first_gauge",
			"type": "gauge",
			"value": 43.44
	    }
	]
	`)
}

func singleElementBatch() io.Reader {
	return bytes.NewBufferString(`
	[
		{
			"id":"single_element",
			"type": "counter",
			"delta": 3
	    }
	]
	`)
}

func (s *UpdatesTestSuite) Do(req *http.Request) *http.Response {
	recorder := httptest.NewRecorder()
	s.handler.ServeHTTP(recorder, req)
	return recorder.Result()
}
