package middleware

import (
	"github.com/mingalevme/feedbacker/infrastructure/log"
	"net/http"
)

func NewLoggingMiddleware(logger log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			logger.Infof("Handling request: %s %s", r.Method, r.URL)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(handler)
	}
}
