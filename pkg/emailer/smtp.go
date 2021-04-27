package emailer

import (
	"fmt"
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/mingalevme/feedbacker/pkg/strutils"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
	"net"
	"strings"
)

type SmtpEmailSender struct {
	host     string  // localhost
	port     uint16  // 25
	username *string // nil
	password *string // nil
	dialer   *gomail.Dialer
	logger   log.Logger
}

func NewSmtpEmailSender(host string, port uint16, username *string, password *string, logger log.Logger) *SmtpEmailSender {
	sender := &SmtpEmailSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		logger:   log.NewNullLogger(),
	}
	if logger != nil {
		sender.logger = logger
	}
	return sender
}

func (s *SmtpEmailSender) Name() string {
	return "smtp"
}

func (s *SmtpEmailSender) Send(from string, to string, subject string, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", strings.Split(to, ",")...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)
	if err := s.smtpDialer().DialAndSend(m); err != nil {
		return errors.Wrap(err, "EmailSender[smtp]: gomail: error while dialing and sending")
	}
	return nil
}

func (s *SmtpEmailSender) smtpDialer() *gomail.Dialer {
	if s.dialer != nil {
		return s.dialer
	}
	s.dialer = &gomail.Dialer{
		Host: s.host,
		Port: int(s.port),
	}
	if strutils.IsNonEmptyString(s.username) {
		s.dialer.Username = *s.username
	}
	if strutils.IsNonEmptyString(s.password) {
		s.dialer.Password = *s.password
	}
	return s.dialer
}

func (s *SmtpEmailSender) Health() error {
	c, err := s.smtpDialer().Dial()
	if err != nil {
		return err
	}
	defer func(c gomail.SendCloser) {
		_ = c.Close()
	}(c)
	return nil
}

func (s *SmtpEmailSender) GetAddr() string {
	return net.JoinHostPort(s.host, fmt.Sprintf("%d", s.port))
}