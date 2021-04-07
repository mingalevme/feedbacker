package notifier

import (
	"encoding/json"
	"fmt"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/pkg/errors"
	"strings"
)

type FeedbackLeftNotifier interface {
	Notify(feedback model.Feedback) error
}

var indent = " "

func feedbackToMessage(f model.Feedback, indent *string) string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%s\n\nContext:\n", f.Text))
	var (
		context []byte
		err error
	)
	if indent != nil {
		context, err = json.MarshalIndent(f, "", *indent)
	} else {
		context, err = json.Marshal(f)
	}
	if err != nil {
		context = []byte(errors.Errorf("Error while marshalling feedback model to json: %v", err).Error())
	}
	b.WriteString(string(context))
	return b.String()
}
