package router

import (
	"github.com/gorilla/mux"
	"github.com/mingalevme/feedbacker/_ddd/infrastructure/di"
	"github.com/mingalevme/feedbacker/_ddd/infrastructure/http/handler"
	"github.com/mingalevme/feedbacker/_ddd/infrastructure/http/middleware"
	netHTTP "net/http"
)

func NewRouter(container di.Container) netHTTP.Handler {
	r := mux.NewRouter()
	//
	r.Handle("/ping", netHTTP.HandlerFunc(handler.Ping)).Methods("GET")
	//
	leaveFeedbackHandler := handler.NewLeaveFeedbackHandler(container.GetLeaveFeedbackService(), container.GetLogger())
	r.Handle("/feedback", middleware.NewJsonRequestMiddleware()(leaveFeedbackHandler)).Methods("POST")
	//
	viewFeedbackHandler := handler.NewViewFeedbackHandler(container.GetViewFeedbackService(), container.GetLogger())
	r.Handle("/feedback/{id}", viewFeedbackHandler).Methods("GET")
	//
	var h netHTTP.Handler
	h = middleware.NewLoggingMiddleware(container.GetLogger())(r)
	h = middleware.NewRecoveryMiddleware(container.GetLogger())(h)
	return h
}
