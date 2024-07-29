package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

var Sign = func(logger *zap.Logger, key string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key == "" {
				h.ServeHTTP(w, r)
				return
			}

			hashEncoded := r.Header.Get("HashSHA256")
			hash, err := base64.StdEncoding.DecodeString(hashEncoded)
			if err != nil {
				err := fmt.Errorf("unable to decode HashSHA256: %w", err)
				logger.Error(err.Error(), zap.Error(err))
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			reqBody, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			bKey := []byte(key)
			bSign := sign(reqBody, bKey)
			if !hmac.Equal(bSign, hash) {
				err := fmt.Errorf("HashSHA256 verification failed")
				logger.Error(err.Error(), zap.Error(err))
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			spy := &ResponseSpy{w: w}
			h.ServeHTTP(spy, r)
			w.Header().Set("HashSHA256", encode(sign(spy.body.Bytes(), bKey)))
			w.WriteHeader(spy.statusCode)
			_, _ = w.Write(spy.body.Bytes())

		})
	}
}

func sign(source []byte, key []byte) []byte {
	hm := hmac.New(sha256.New, key)
	hm.Write(source)
	return hm.Sum(nil)
}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
