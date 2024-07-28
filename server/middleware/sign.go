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

			hm := hmac.New(sha256.New, []byte(key))
			hm.Write(reqBody)
			sign := hm.Sum(nil)
			if !hmac.Equal(sign, hash) {
				err := fmt.Errorf("HashSHA256 verification failed")
				logger.Error(err.Error(), zap.Error(err))
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
