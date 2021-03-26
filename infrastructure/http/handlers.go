package http

import (
	"github.com/mingalevme/feedbacker/app/services"
	netHTTP "net/http"
)

type Handler netHTTP.Handler

type SimpleHandler func(w netHTTP.ResponseWriter, r *netHTTP.Request)

type SimpleRequestHandler interface {
	ServeHTTP(w netHTTP.ResponseWriter, r *netHTTP.Request)
}

type requestHandler struct {
	handler SimpleHandler
}

func NewRequestHandler(handler SimpleHandler) SimpleRequestHandler {
	return &requestHandler{
		handler: handler,
	}
}

func (h requestHandler) ServeHTTP(w netHTTP.ResponseWriter, r *netHTTP.Request) {
	h.handler(w, r)
}

// ---------------------------------------------------------------------------------------------------------------------

func ping(w netHTTP.ResponseWriter, r *netHTTP.Request) {
	_, _ = w.Write([]byte("Pong"))
}

// ---------------------------------------------------------------------------------------------------------------------

type LeaveFeedbackHandler Handler

type leaveFeedbackHandler struct {
	service services.LeaveFeedbackService
}

func NewLeaveFeedbackHandler(service services.LeaveFeedbackService) LeaveFeedbackHandler {
	return &leaveFeedbackHandler{
		service: service,
	}
}

func (handler leaveFeedbackHandler) ServeHTTP(w netHTTP.ResponseWriter, r *netHTTP.Request) {
	//body := make([]byte, env.MaxPostRequestBodyLength()+1)
	//n, err := io.ReadFull(r.Body, body)
	//if err != nil {
	//	if err != io.ErrUnexpectedEOF {
	//		return errors.Wrap(err, "Error while reading request body")
	//	}
	//	body = body[:n]
	//}
	//if n < 1 {
	//	return makeBadRequestHttpError("Request body is empty", nil)
	//}
	//if uint(n) > env.MaxPostRequestBodyLength() {
	//	return makeRequestEntityTooLargeHttpError("Request body is too large", nil)
	//}
}

// ---------------------------------------------------------------------------------------------------------------------
