package notifier

import "github.com/mingalevme/feedbacker/internal/app/model"

type ErrNotifier struct {
	err error
}

func (s *ErrNotifier) Name() string {
	return "error"
}

func (s *ErrNotifier) Health() error {
	return s.err
}

func (s *ErrNotifier) Notify(f model.Feedback) error {
	return s.err
}

func NewErrNotifier(err error) *ErrNotifier {
	return &ErrNotifier{
		err: err,
	}
}
