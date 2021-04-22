// +build testing

package test

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/internal/app/repository"
)

type HealthErrorRepository struct {
	Err error
}

func (s *HealthErrorRepository) Name() string {
	return "error"
}

func (s *HealthErrorRepository) Add(data repository.AddFeedbackData) (model.Feedback, error) {
	panic("implement me")
}

func (s *HealthErrorRepository) GetById(id int) (model.Feedback, error) {
	panic("implement me")
}

func (s *HealthErrorRepository) Health() error {
	return s.Err
}
