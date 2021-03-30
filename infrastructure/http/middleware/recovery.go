package middleware

import (
	"fmt"
	"github.com/mingalevme/feedbacker/infrastructure/config"
	"github.com/mingalevme/feedbacker/infrastructure/log"
	"net/http"
	"runtime/debug"
)

func NewRecoveryMiddleware(logger log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					//var err error
					//switch v := r.(type) {
					//case error:
					//	err = v
					//case string:
					//	err = errors.New(v)
					//default:
					//	err = errors.New(fmt.Sprintf("Unkown panic error (%T): %+v", v, v))
					//}
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					if config.GetInstance().IsDebug() {
						http.Error(w, fmt.Sprintf("Panic (%T): %+v\nStacktrace:\n%s", r, r, debug.Stack()), http.StatusInternalServerError)
					} else {
						logger.Errorf("Panic (%T): %+v\nStacktrace:\n%s", r, r, debug.Stack())
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					}
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(handler)
	}
}
