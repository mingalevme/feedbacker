package notifier

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type SlackPanicHTTPClient struct{}

func (s *SlackPanicHTTPClient) Do(req *http.Request) (*http.Response, error) {
	panic(req)
}

func TestSlackSendNoError(t *testing.T) {
	client := slack.New(
		"xoxb-123456789012-1234567890123-1234567890qwertyuiopasdf",
		slack.OptionHTTPClient(&SlackPanicHTTPClient{}),
	)
	notifier := NewSlackNotifier(client, "QWERTYUIO")
	f := model.MakeFeedback()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Message has not been sent to Slack")
		}
		req, ok := r.(*http.Request)
		if !ok {
			panic(r)
		}
		assert.Equal(t, "https://slack.com/api/chat.postMessage", req.URL.String())
	}()
	_ = notifier.Notify(f)
}
