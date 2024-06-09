package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiveness(t *testing.T) {
	h := NewHandler()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/liveness", nil)

	h.ServeHTTP(recorder, req)
	resp := recorder.Result()
	defer closeBody(resp)

	b, _ := io.ReadAll(recorder.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "alive", string(b))
}

func closeBody(resp *http.Response) {
	_ = resp.Body.Close()
}
