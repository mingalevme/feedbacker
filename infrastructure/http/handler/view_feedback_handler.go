package handler

import (
	"github.com/mingalevme/feedbacker/app/service"
	"github.com/mingalevme/feedbacker/infrastructure/http/responder"
	"github.com/mingalevme/feedbacker/infrastructure/log"
	"github.com/pkg/errors"
	"net/http"
)

type ViewFeedbackHandler Handler

type viewFeedbackHandler struct {
	service service.ViewFeedbackService
	logger log.Logger
}

func NewViewFeedbackHandler(service service.ViewFeedbackService, logger log.Logger) ViewFeedbackHandler {
	return viewFeedbackHandler{
		service: service,
		logger: logger,
	}
}

func (h viewFeedbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := service.ViewFeedbackData{}
	if err := requestToData(r, &data); err != nil {
		h.logger.Error("Error while deserialize request to params", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	f, err := h.service.Handle(data)
	if errors.Is(err, service.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if errors.Is(err, service.ErrUnprocessableEntity) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, errors.Wrap(err, "Error while leaving feedback").Error(), http.StatusBadRequest)
		return
	}
	resp := responder.NewJSONResponder(h.logger)
	_ = resp.Respond(w, f)
}
