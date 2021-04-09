package notifier

import (
	"fmt"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/pkg/errors"
	"strings"
)

type StackNotifier struct {
	notifiers []Notifier
}

func (s *StackNotifier) Name() string {
	var names []string
	for _, notifier := range s.notifiers {
		names = append(names, notifier.Name())
	}
	return fmt.Sprintf("stack: %s", strings.Join(names, ", "))
}

func (s *StackNotifier) Health() error {
	var result []string
	for _, notifier := range s.notifiers {
		if err := notifier.Health(); err != nil {
			result = append(result, fmt.Sprintf("%s: %s", notifier.Name(), err.Error()))
		}
	}
	if len(result) > 0 {
		return errors.New(strings.Join(result, ";"))
	}
	return nil
}

func (s *StackNotifier) Notify(f model.Feedback) error {
	var result []string
	for _, notifier := range s.notifiers {
		if err := notifier.Notify(f); err != nil {
			result = append(result, fmt.Sprintf("%s: %s", notifier.Name(), err.Error()))
		}
	}
	if len(result) > 0 {
		return errors.New(strings.Join(result, "; "))
	}
	return nil
}

func (s *StackNotifier) Add(n Notifier) {
	s.notifiers = append(s.notifiers, n)
}

func NewStackNotifier() *StackNotifier {
	return &StackNotifier{
		notifiers: []Notifier{},
	}
}
