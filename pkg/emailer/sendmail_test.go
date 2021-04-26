package emailer

import (
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

// todo: mock command executing
func TestSendmailSending(t *testing.T) {
	l := log.NewArrayLogger(log.LevelDebug)
	//l := log.NewStdoutLogger(log.LevelDebug)
	//cmd := "/bin/cat"
	cmd := "/usr/sbin/sendmail"
	sender := NewSendmailEmailSender(cmd, l)
	err := sender.Send("from@example.com", "to@example.com", "Test Subject", "Test body")
	assert.NoError(t, err)
}
