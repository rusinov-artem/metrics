package middleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

var Decrypt = func(privateKey *rsa.PrivateKey, logger *zap.Logger) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if privateKey == nil {
				h.ServeHTTP(w, r)
			}

			body, _ := io.ReadAll(r.Body)
			r.Body.Close()
			decryptedBody, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, body)
			if err != nil {
				logger.Error(fmt.Sprintf("unable to decrypt message: %v", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))
			h.ServeHTTP(w, r)
		})
	}
}
