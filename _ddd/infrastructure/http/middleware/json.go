package middleware

import (
	"github.com/mingalevme/feedbacker/_ddd/infrastructure/util"
	"net/http"
)

func NewJsonRequestMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			h := util.ParseHeader(r.Header.Get("Content-Type"))
			v, ok := h["application/json"]
			// If header does not contain "application/json" or contains "application/json=XXX"
			if ok != true || v != nil {
				http.Error(w, "Unknown media type", http.StatusUnsupportedMediaType)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(handler)
	}
}
