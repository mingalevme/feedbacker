package task

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/internal/app/service/notifier"
	"github.com/mingalevme/feedbacker/pkg/dispatcher"
)

func NewNotifyTask(notifier notifier.Notifier, f model.Feedback) dispatcher.Task {
	return func() error {
		return notifier.Notify(f)
	}
}
