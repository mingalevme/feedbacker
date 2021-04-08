package emailer

import (
	"fmt"
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/pkg/errors"
)

type arrayEmailItem struct {
	from string
	to string
	message string
}

type ArrayEmailSender struct {
	storage []arrayEmailItem
	logger  log.Logger
}

func NewArrayEmailSender(logger log.Logger) *ArrayEmailSender {
	if logger == nil {
		panic(errors.New("logger is empty"))
	}
	sender := &ArrayEmailSender{
		storage: []arrayEmailItem{},
		logger:   logger,
	}
	return sender
}

func (s *ArrayEmailSender) Send(from string, to string, subject string, message string) error {
	context := map[string]interface{}{
		"_sender": fmt.Sprintf("%T", s),
		"from": from,
		"to": to,
		"message": message,
	}
	s.logger.WithFields(context).Info("Sending email")
	s.storage = append(s.storage, arrayEmailItem{
		from:    from,
		to:      to,
		message: message,
	})
	s.logger.WithFields(context).Info("Email has been sent successfully")
	return nil

}