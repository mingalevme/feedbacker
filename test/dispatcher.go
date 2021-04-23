// +build testing

package test

import (
	"github.com/mingalevme/feedbacker/pkg/dispatcher"
)

type HealthErrorDispatcher struct {
	Err error
}

func (s *HealthErrorDispatcher) Name() string {
	return "error"
}

func (s *HealthErrorDispatcher) Enqueue(t dispatcher.Task) error {
	return nil
}

func (s *HealthErrorDispatcher) Run() error {
	return nil
}

func (s *HealthErrorDispatcher) Stop() error {
	return nil
}

func (s *HealthErrorDispatcher) Health() error {
	return s.Err
}
