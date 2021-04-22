// +build testing

package test

import "github.com/mingalevme/feedbacker/internal/app/model"

type HealthErrorNotifier struct {
	Err error
}

func (s *HealthErrorNotifier) Name() string {
	return "error"
}

func (s *HealthErrorNotifier) Health() error {
	return s.Err
}

func (s *HealthErrorNotifier) Notify(model.Feedback) error {
	panic("implement me")
}
