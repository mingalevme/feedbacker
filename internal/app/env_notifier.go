package app

import (
	"github.com/mingalevme/feedbacker/internal/app/service/notifier"
	"github.com/pkg/errors"
	"strings"
)

func (s *Container) Notifier() notifier.Notifier {
	if s.notifier != nil {
		return s.notifier
	}
	channel := s.EnvVarBag.Get("NOTIFIER_CHANNEL", "email")
	s.notifier = s.newNotifierChannel(channel)
	return s.notifier
}

func (s *Container) newNotifierChannel(channel string) notifier.Notifier {
	if channel == "email" {
		return notifier.NewEmailNotifier(s.EmailSender(), s.NotifierEmailFrom(), s.NotifierEmailTo(), s.NotifierEmailSubjectTemplate(), s.Logger())
	}
	if channel == "slack" {
		channelID := s.EnvVarBag.Require("NOTIFIER_SLACK_CHANNEL_ID")
		return notifier.NewSlackNotifier(s.Slack(), channelID)
	}
	if channel == "array" {
		return notifier.NewArrayNotifier(s.Logger())
	}
	if channel == "null" {
		return notifier.NewNullNotifier()
	}
	if channel == "stack" {
		return s.newStackNotifier()
	}
	panic(errors.Errorf("Unsupported notifier channel: %s", channel))
}

func (s *Container) newStackNotifier() *notifier.StackNotifier {
	n := notifier.NewStackNotifier()
	channels := strings.Split(s.EnvVarBag.Require("NOTIFIER_STACK_CHANNELS"), ",")
	for _, channel := range channels {
		channel = strings.TrimSpace(channel)
		if channel == "" {
			continue
		}
		if channel == "stack" {
			panic(errors.Errorf("stack channel recursion"))
		}
		n.Add(s.newNotifierChannel(channel))
	}
	return n
}

func (s *Container) NotifierEmailSubjectTemplate() string {
	return s.EnvVarBag.Get("NOTIFIER_EMAIL_SUBJECT_TEMPLATE", "Feedback %{InstallationID}s")
}
