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

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/liveness", nil)

	h.ServeHTTP(resp, req)

	b, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Result().StatusCode)
	assert.Equal(t, "alive", string(b))
}
