package notifier

import (
	"fmt"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/internal/app/service/emailer"
	"github.com/mingalevme/feedbacker/internal/app/service/log"
	util2 "github.com/mingalevme/feedbacker/pkg/util"
	"github.com/pkg/errors"
	"strings"
)

type EmailFeedbackLeftNotifier struct {
	sender  emailer.EmailSender
	from    string
	to      string
	subject string
	logger  log.Logger
}

// Sync
func (s *EmailFeedbackLeftNotifier) Notify(f model.Feedback) error {
	s.logger.WithField("_notifier", fmt.Sprintf("%T", s)).Infof("Notifying:\n%s", feedbackToMessage(f, nil))
	replacement := map[string]interface{}{
		"InstallationID": "",
	}
	if f.Customer != nil && f.Customer.InstallationID != nil {
		replacement["InstallationID"] = *f.Customer.InstallationID
	}
	subject := strings.TrimSpace(util2.Sprintf(s.subject, replacement))
	return s.sender.Send(s.from, s.to, subject, feedbackToMessage(f, &indent))
}

func NewEmailFeedbackLeftNotifier(sender emailer.EmailSender, from string, to string, subject string, logger log.Logger) *EmailFeedbackLeftNotifier {
	if sender == nil {
		panic(errors.New("`sender` is nil"))
	}
	if util2.IsEmptyString(from) {
		panic(errors.New("`from` is empty"))
	}
	if util2.IsEmptyString(to) {
		panic(errors.New("`to` is empty"))
	}
	if util2.IsEmptyString(subject) {
		panic(errors.New("`subject` is empty"))
	}
	if logger == nil {
		panic(errors.New("`logger` is nil"))
	}
	notifier := &EmailFeedbackLeftNotifier{
		sender:  sender,
		subject: subject,
		from:    from,
		to:      to,
		logger:  logger,
	}
	return notifier
}
