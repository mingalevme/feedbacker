package emailer

import (
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type sendmailCommandMock struct {
	isStarted bool
	stdin     *sendmailCommandStdinMock
}

func (s *sendmailCommandMock) Start() error {
	if s.isStarted {
		return errors.New("Command already has been started")
	}
	s.isStarted = true
	return nil
}

func (s *sendmailCommandMock) StdinPipe() (io.WriteCloser, error) {
	if s.stdin != nil {
		return nil, errors.New("Stdin already set")
	}
	if s.isStarted {
		return nil, errors.New("StdinPipe after command started")
	}
	s.stdin = newSendmailCommandStdinMock()
	return s.stdin, nil
}

func (s *sendmailCommandMock) Wait() error {
	return nil
}

type sendmailCommandStdinMock struct {
	isClosed bool
	storage  string
}

func newSendmailCommandStdinMock() *sendmailCommandStdinMock {
	return &sendmailCommandStdinMock{
		storage: "",
	}
}

func (s *sendmailCommandStdinMock) Write(p []byte) (n int, err error) {
	if s.isClosed {
		return 0, errors.New("Pipe is closed")
	}
	s.storage = s.storage + string(p)
	return len(p), nil
}

func (s *sendmailCommandStdinMock) Close() error {
	s.isClosed = true
	return nil
}

func TestSendmailSendingNoError(t *testing.T) {
	l := log.NewArrayLogger(log.LevelDebug)
	cmd := "/usr/sbin/sendmail"
	sender := NewSendmailEmailSender(cmd, l)
	mock := &sendmailCommandMock{}
	sender.factory = func(name string, arg ...string) SendmailCommand {
		return mock
	}
	err := sender.Send("from@example.com", "to@example.com", "Test Subject", "Test body")
	assert.Len(t, l.Storage(), 0)
	assert.NoError(t, err)
	assert.True(t, mock.isStarted)
	assert.True(t, mock.stdin.isClosed)
	assert.Contains(t, mock.stdin.storage, "From: from@example.com")
	assert.Contains(t, mock.stdin.storage, "To: to@example.com")
	assert.Contains(t, mock.stdin.storage, "Subject: Test Subject")
	assert.Contains(t, mock.stdin.storage, "\r\n\r\nTest body")
}
