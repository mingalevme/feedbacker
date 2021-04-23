// +build testing

package app

import (
	"github.com/mingalevme/feedbacker/internal/app/repository"
	"github.com/mingalevme/feedbacker/internal/app/service/notifier"
	"github.com/mingalevme/feedbacker/pkg/dispatcher"
)

func (s *Container) SetFeedbackRepository(r repository.Feedback) {
	s.repository = r
}

func (s *Container) SetNotifier(n notifier.Notifier) {
	s.notifier = n
}

func (s *Container) SetDispatcher(d dispatcher.Dispatcher) {
	s.dispatcher = d
}
