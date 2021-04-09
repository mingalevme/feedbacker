package notifier

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/pkg/emailer"
	"github.com/mingalevme/feedbacker/pkg/strutils"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestEmailSend(t *testing.T) {
	sender := emailer.NewArrayEmailSender()
	notifier := NewEmailNotifier(sender, "from@mail.com", "to.mail.com", "  New %{InstallationID}s", nil)
	f := model.MakeFeedback()
	f.Customer.InstallationID = strutils.StrToPointerStr("FooBar")
	err := notifier.Notify(f)
	assert.NoError(t, err)
	assert.Len(t, sender.Storage, 1)
	assert.Equal(t, "New FooBar", sender.Storage[0].Subject)
}

func TestEmailHealthNoErr(t *testing.T) {
	sender := emailer.NewArrayEmailSender()
	notifier := NewEmailNotifier(sender, "from@mail.com", "to.mail.com", "  New %{InstallationID}s", nil)
	err := notifier.Health()
	assert.NoError(t, err)
}

func TestEmailHealthErr(t *testing.T) {
	sender := emailer.NewSmtpEmailSender("localhost", uint16(49152), nil,  nil, nil)
	notifier := NewEmailNotifier(sender, "from@mail.com", "to.mail.com", "New %{InstallationID}s", nil)
	err := notifier.Health()
	var e *net.OpError
	assert.ErrorAs(t, err, &e)
}
