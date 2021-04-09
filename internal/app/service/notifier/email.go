package notifier

import (
	"fmt"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/pkg/emailer"
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/mingalevme/feedbacker/pkg/util"
	"github.com/pkg/errors"
	"strings"
)

type EmailNotifier struct {
	sender  emailer.EmailSender
	from    string
	to      string
	subject string
	logger  log.Logger
}

func (s *EmailNotifier) Name() string {
	return "emailer"
}

func (s *EmailNotifier) Health() error {
	return s.sender.Health()
}

// Sync
func (s *EmailNotifier) Notify(f model.Feedback) error {
	s.logger.WithField("_notifier", fmt.Sprintf("%T", s)).Infof("Notifying:\n%s", feedbackToMessage(f, nil))
	replacement := map[string]interface{}{
		"InstallationID": "",
	}
	if f.Customer != nil && f.Customer.InstallationID != nil {
		replacement["InstallationID"] = *f.Customer.InstallationID
	}
	subject := strings.TrimSpace(util.Sprintf(s.subject, replacement))
	return s.sender.Send(s.from, s.to, subject, feedbackToMessage(f, &indent))
}

func NewEmailNotifier(sender emailer.EmailSender, from string, to string, subject string, logger log.Logger) *EmailNotifier {
	if sender == nil {
		panic(errors.New("`sender` is nil"))
	}
	if util.IsEmptyString(from) {
		panic(errors.New("`from` is empty"))
	}
	if util.IsEmptyString(to) {
		panic(errors.New("`to` is empty"))
	}
	if util.IsEmptyString(subject) {
		panic(errors.New("`subject` is empty"))
	}
	notifier := &EmailNotifier{
		sender:  sender,
		subject: subject,
		from:    from,
		to:      to,
		logger:  log.NewNullLogger(),
	}
	if logger != nil {
		notifier.logger = logger
	}
	return notifier
}
