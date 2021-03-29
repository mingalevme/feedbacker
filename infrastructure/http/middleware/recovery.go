package middleware

import (
	"fmt"
	"github.com/mingalevme/feedbacker/infrastructure/log"
	"net/http"
)

func NewRecoveryMiddleware(logger log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					logger.Panic(err)
					http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(handler)
	}
}
