package notifier

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/pkg/emailer"
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/mingalevme/feedbacker/pkg/strutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendingEmail(t *testing.T) {
	logger := log.NewNullLogger()
	sender := emailer.NewArrayEmailSender()
	notifier := NewEmailNotifier(sender, "from@mail.com", "to.mail.com", "  New %{InstallationID}s", logger)
	f := model.MakeFeedback()
	f.Customer.InstallationID = strutils.StrToPointerStr("FooBar")
	err := notifier.Notify(f)
	assert.NoError(t, err)
	assert.Len(t, sender.Storage, 1)
	assert.Equal(t, "New FooBar", sender.Storage[0].Subject)
}

