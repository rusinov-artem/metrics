package middleware

import (
	"bytes"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ResponseSpy struct {
	w          http.ResponseWriter
	body       bytes.Buffer
	statusCode int
}

func (r *ResponseSpy) Header() http.Header {
	return r.w.Header()
}

func (r *ResponseSpy) Write(bytes []byte) (int, error) {
	r.body.Write(bytes)
	return r.w.Write(bytes)
}

func (r *ResponseSpy) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.w.WriteHeader(statusCode)
}

var Logger = func(logger *zap.Logger) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			spy := &ResponseSpy{w: w}
			h.ServeHTTP(spy, r)
			dur := time.Since(start)
			logger.Info("handling request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Duration("duration", dur),
				zap.Int("statusCode", spy.statusCode),
				zap.Int("size", spy.body.Len()),
			)
		})
	}
}
