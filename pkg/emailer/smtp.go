package emailer

import (
	"fmt"
	"github.com/mingalevme/feedbacker/pkg/log"
	util2 "github.com/mingalevme/feedbacker/pkg/util"
	"github.com/pkg/errors"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type SmtpEmailSender struct {
	host     string  // localhost
	port     uint16  // 25
	username *string // nil
	password *string // nil
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

func (s *SmtpEmailSender) Health() error {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", s.GetAddr(), timeout)
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	return err
}

func (s *SmtpEmailSender) Send(from string, to string, subject string, message string) error {
	context := map[string]interface{}{
		"_sender": fmt.Sprintf("%T", s),
		"from":    from,
		"to":      to,
		"message": message,
	}
	s.logger.WithFields(context).Info("Sending email")
	subject = strings.ReplaceAll(subject, "\n", " ")
	message = fmt.Sprintf("Subject: %s\n\n%s", subject, message)
	err := smtp.SendMail(
		s.GetAddr(),
		s.getAuth(),
		from,
		strings.Split(to, ","),
		[]byte(message),
	)
	if err != nil {
		return wrapSmtpError(errors.Wrap(err, "sending mail"))
	}
	s.logger.WithFields(context).Info("Email has been sent successfully")
	return nil
	//c, err := smtp.Dial(s.GetAddr())
	//if err != nil {
	//	return wrapSmtpError(errors.Wrap(err, "opening smtp connection"))
	//}
	//defer func() {
	//	_ = c.Close()
	//}()
	//if err = c.Mailfrom); err != nil {
	//	return wrapSmtpError(errors.Wrap(err, "setting `From`"))
	//}
	//if err = c.Rcpt(to); err != nil {
	//	return wrapSmtpError(errors.Wrap(err, "setting `To`"))
	//}
	//body, err := c.Data()
	//if err != nil {
	//	return wrapSmtpError(errors.Wrap(err, "initializing writer"))
	//}
	//defer func() {
	//	_ = body.Close()
	//}()
	//buf := bytes.NewBufferString(message)
	//if _, err = buf.WriteTo(body); err != nil {
	//	return wrapSmtpError(errors.Wrap(err, "writing body"))
	//}
	//return nil
}

func (s *SmtpEmailSender) getAuth() smtp.Auth {
	if util2.IsEmptyString(s.username) {
		return nil
	}
	var username = *s.username
	var password = ""
	if util2.IsNonEmptyString(s.password) {
		password = *s.password
	}
	return smtp.PlainAuth(
		"",
		username,
		password,
		s.host,
	)
}

func (s *SmtpEmailSender) GetAddr() string {
	return net.JoinHostPort(s.host, fmt.Sprintf("%d", s.port))
	//return fmt.Sprintf("%s:%d", s.host, s.port)
}

func wrapSmtpError(err error) error {
	return errors.Wrap(err, "smtp email sender")
}
