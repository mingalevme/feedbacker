package handler

import (
	"encoding/json"
	"github.com/mingalevme/feedbacker/app/service"
	"github.com/mingalevme/feedbacker/infrastructure/log"
	"github.com/pkg/errors"
	"net/http"
)

type LeaveFeedbackHandler Handler

type leaveFeedbackHandler struct {
	service service.LeaveFeedbackService
	logger log.Logger
}

func NewLeaveFeedbackHandler(service service.LeaveFeedbackService, logger log.Logger) LeaveFeedbackHandler {
	return &leaveFeedbackHandler{
		service: service,
		logger: logger,
	}
}

func (handler leaveFeedbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := service.LeaveFeedbackData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	feedback, err := handler.service.Handle(data)
	if err != nil {
		http.Error(w, errors.Wrap(err, "Error while leaving feedback").Error(), http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(feedback); err != nil {
		http.Error(w, errors.Wrap(err, "Error while leaving feedback").Error(), http.StatusBadRequest)
		return
	}
}
