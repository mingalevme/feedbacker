package emailer

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestSmtpName(t *testing.T) {
	smtpEmailSender := NewSmtpEmailSender("127.0.0.1", uint16(0), nil, nil, nil)
	assert.Equal(t, "smtp", smtpEmailSender.Name())
}

func TestSmtpHealthNetOpError(t *testing.T) {
	smtpEmailSender := NewSmtpEmailSender("127.0.0.1", uint16(0), nil, nil, nil)
	err := smtpEmailSender.Health()
	var e *net.OpError
	assert.ErrorAs(t, err, &e)
}

func TestSmtpHealthNoError(t *testing.T) {
	l, err := net.Listen("tcp4", "127.0.0.1:49152")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	go func() {
		conn, err := l.Accept()
		assert.NoError(t, err)
		defer conn.Close()
		_, err = conn.Write([]byte("220 Ok"))
	}()
	smtpEmailSender := NewSmtpEmailSender("127.0.0.1", uint16(49152), nil, nil, nil)
	err = smtpEmailSender.Health()
	assert.NoError(t, err)
}
