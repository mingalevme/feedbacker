package repository

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("model is not found")
var ErrUnprocessableEntity = errors.New("unprocessable entity")

type AddFeedbackData struct {
	model.Feedback
}

func (s AddFeedbackData) Validate() error {
	if s.Service == "" {
		return errors.Wrap(ErrUnprocessableEntity, "Service is empty")
	}
	if s.Text == "" {
		return errors.Wrap(ErrUnprocessableEntity, "Text is empty")
	}
	return nil
}

type Feedback interface {
	Name() string
	Add(data AddFeedbackData) (model.Feedback, error)
	GetById(id int) (model.Feedback, error)
	Health() error
}
