// +build testing

package app

import (
	"github.com/mingalevme/feedbacker/internal/app/repository"
	"github.com/mingalevme/feedbacker/internal/app/service/notifier"
)

func (s *Container) SetFeedbackRepository(r repository.Feedback) {
	s.repository = r
}

func (s *Container) SetNotifier(n notifier.Notifier) {
	s.notifier = n
}
