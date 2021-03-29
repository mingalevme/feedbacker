package router

import (
	"github.com/gorilla/mux"
	"github.com/mingalevme/feedbacker/infrastructure/di"
	"github.com/mingalevme/feedbacker/infrastructure/http/handler"
	"github.com/mingalevme/feedbacker/infrastructure/http/middleware"
	netHTTP "net/http"
)

func NewRouter(container di.Container) netHTTP.Handler {
	r := mux.NewRouter()
	//
	leaveFeedbackHandler := handler.NewLeaveFeedbackHandler(container.GetLeaveFeedbackService(), container.GetLogger())
	r.Handle("/feedback", middleware.NewJsonRequestMiddleware()(leaveFeedbackHandler)).Methods("POST")
	//
	r.Handle("/ping", netHTTP.HandlerFunc(handler.Ping)).Methods("GET")
	//
	var h netHTTP.Handler
	h = middleware.NewRecoveryMiddleware(container.GetLogger())(r)
	h = middleware.NewLoggingMiddleware(container.GetLogger())(h)
	return h
}
