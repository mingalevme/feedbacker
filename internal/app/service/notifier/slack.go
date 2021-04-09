package notifier

import (
	"fmt"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	client    *slack.Client
	channelID string
}

func (s *SlackNotifier) Name() string {
	return "slack"
}

func (s *SlackNotifier) Health() error {
	return nil
}

func (s *SlackNotifier) Notify(f model.Feedback) error {
	message := feedbackToMessage(f, &indent)
	channelID, timestamp, err := s.client.PostMessage(
		s.channelID,
		slack.MsgOptionText(message, false),
	)
	fmt.Printf("Slack: '%v' '%v' '%v'", channelID, timestamp, err)
	return err
}

func NewSlackNotifier(client *slack.Client, channel string) *SlackNotifier {
	return &SlackNotifier{
		client:    client,
		channelID: channel,
	}
}
