package view_feedback

import (
	"github.com/mingalevme/feedbacker/internal/app"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/internal/app/repository"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New(repository.ErrNotFound.Error())

type ViewFeedback struct {
	env app.Env
}

func New(env app.Env) *ViewFeedback {
	return &ViewFeedback{
		env: env,
	}
}

func (s *ViewFeedback) ViewFeedback(id int) (model.Feedback, error) {
	f, err := s.env.FeedbackRepository().GetById(id)
	if err == repository.ErrNotFound {
		return f, ErrNotFound
	}
	return f, err
}
