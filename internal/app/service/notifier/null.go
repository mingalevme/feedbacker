package notifier

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
)

type NullNotifier struct {
}

func (s *NullNotifier) Name() string {
	return "null"
}

func (s *NullNotifier) Health() error {
	return nil
}

func (s *NullNotifier) Notify(f model.Feedback) error {
	return nil
}

func NewNullNotifier() *NullNotifier {
	return &NullNotifier{}
}
