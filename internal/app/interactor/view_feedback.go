package interactor

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/internal/app/repository"
)

func (s *Interactor) ViewFeedback(id int) (model.Feedback, error) {
	f, err := s.env.FeedbackRepository().GetById(id)
	if err == repository.ErrNotFound {
		return f, ErrNotFound
	}
	return f, err
}
