package server

import (
	"log"
	"net/http"
)

func NewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		if r.URL.Path == "/liveness" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("alive"))
		}
	})
}
