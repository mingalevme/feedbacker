package handler

import (
	"encoding/json"
	"github.com/mingalevme/feedbacker/_ddd/app/service"
	"github.com/mingalevme/feedbacker/_ddd/infrastructure/http/responder"
	"github.com/mingalevme/feedbacker/internal/app/service/log"
	"github.com/pkg/errors"
	"net/http"
)

type LeaveFeedbackHandler Handler

type leaveFeedbackHandler struct {
	service service.LeaveFeedbackService
	logger  log.Logger
}

func NewLeaveFeedbackHandler(service service.LeaveFeedbackService, logger log.Logger) LeaveFeedbackHandler {
	return leaveFeedbackHandler{
		service: service,
		logger: logger,
	}
}

func (h leaveFeedbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := service.LeaveFeedbackData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	feedback, err := h.service.Handle(data)
	if err != nil {
		http.Error(w, errors.Wrap(err, "Error while leaving feedback").Error(), http.StatusBadRequest)
		return
	}
	resp := responder.NewJSONResponder(h.logger)
	_ = resp.Respond(w, feedback)
}
