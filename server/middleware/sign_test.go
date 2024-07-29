package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/rusinov-artem/metrics/test/logger"
)

type SignTestSuite struct {
	suite.Suite
	logger *zap.Logger
	logs   *bytes.Buffer
}

func Test_SignMiddleware(t *testing.T) {
	suite.Run(t, &SignTestSuite{})
}

func (s *SignTestSuite) SetupTest() {
	s.logger, s.logs = logger.SpyLogger()
}

func (s *SignTestSuite) Test_DoNothingIfKeyIsEmpty() {
	m := Sign(s.logger, "")
	mux := m(http.HandlerFunc((&DummyHandler{}).Handle))

	req := httptest.NewRequest(http.MethodPost, "/target", bytes.NewBufferString(`my request`))
	resp := httptest.NewRecorder()

	mux.ServeHTTP(resp, req)
	defer func() { _ = resp.Result().Body.Close() }()

	s.Equal(200, resp.Code)
	s.Equal("My body", resp.Body.String())
	s.Equal("Value", resp.Header().Get("My-Header"))
}

func (s *SignTestSuite) Test_RejectRequestWithoutSign() {
	m := Sign(s.logger, "key")
	mux := m(http.HandlerFunc((&DummyHandler{}).Handle))

	req := httptest.NewRequest(http.MethodPost, "/target", bytes.NewBufferString(`my request`))
	resp := httptest.NewRecorder()

	mux.ServeHTTP(resp, req)
	defer func() { _ = resp.Result().Body.Close() }()

	s.Equal(400, resp.Code)
	s.Contains(resp.Body.String(), "HashSHA256 verification failed")
	s.Equal("", resp.Header().Get("My-Header"))
}

func (s *SignTestSuite) Test_RejectRequestWithBadlyEncodedSign() {
	m := Sign(s.logger, "key")
	mux := m(http.HandlerFunc((&DummyHandler{}).Handle))

	req := httptest.NewRequest(http.MethodPost, "/target", bytes.NewBufferString(`my request`))
	resp := httptest.NewRecorder()
	req.Header.Set("HashSHA256", "/.dInvalid")

	mux.ServeHTTP(resp, req)
	defer func() { _ = resp.Result().Body.Close() }()

	s.Equal(400, resp.Code)
	s.Contains(resp.Body.String(), "illegal base64 data at input")
	s.Equal("", resp.Header().Get("My-Header"))
}

func (s *SignTestSuite) Test_PassThroughRequestWithCorrectSign() {
	m := Sign(s.logger, "key")
	mux := m(http.HandlerFunc((&DummyHandler{}).Handle))

	reqBody := bytes.NewBufferString(`my request`)
	req := httptest.NewRequest(http.MethodPost, "/target", reqBody)
	req.Header.Set("HashSHA256", encode(sign(reqBody.Bytes(), []byte("key"))))
	resp := httptest.NewRecorder()

	mux.ServeHTTP(resp, req)
	defer func() { _ = resp.Result().Body.Close() }()

	s.Equal(200, resp.Code)
	s.Equal("My body", resp.Body.String())
	s.Equal("Value", resp.Header().Get("My-Header"))
}

type DummyHandler struct {
}

func (d *DummyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("My-Header", "Value")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(`My body`))
}
