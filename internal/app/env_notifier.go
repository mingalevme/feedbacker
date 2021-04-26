package app

import (
	"fmt"
	"github.com/mingalevme/feedbacker/internal/app/service/notifier"
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/mingalevme/feedbacker/pkg/util"
	"github.com/pkg/errors"
	"os"
	"os/user"
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
	if channel == "log" {
		level, err := log.ParseLevel(s.EnvVarBag.Get("NOTIFIER_LOG_LEVEL", "debug"))
		if  err != nil {
			panic(errors.Errorf("invalid NOTIFIER_LOG_LEVEL env-var"))
		}
		return notifier.NewLogNotifier(s.Logger(), level)
	}
	if channel == "array" {
		return notifier.NewArrayNotifier()
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

func (s *Container) NotifierEmailFrom() string {
	def := func() string {
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		h, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("%s@%s", u.Username, h)
	}
	from := s.EnvVarBag.Get("NOTIFIER_EMAIL_FROM", "")
	if util.IsEmptyString(from) {
		return def()
	}
	return from
}

func (s *Container) NotifierEmailTo() string {
	to := s.EnvVarBag.Require("NOTIFIER_EMAIL_TO")
	return to
}

func (s *Container) NotifierEmailSubjectTemplate() string {
	return s.EnvVarBag.Get("NOTIFIER_EMAIL_SUBJECT_TEMPLATE", "Feedback %{InstallationID}s")
}
