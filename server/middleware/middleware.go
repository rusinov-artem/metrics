package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func (t *Middleware) Wrap(m Middleware) Middleware {
	parent := *t
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m((parent)(h)).ServeHTTP(w, r)
		})
	}
}
