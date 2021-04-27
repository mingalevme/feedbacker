package emailer

import (
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
	"io"
	"os"
	"os/exec"
)

type SendmailEmailSender struct {
	cmd     string
	logger  log.Logger
	factory SendmailCommandFactory
}

func NewSendmailEmailSender(cmd string, logger log.Logger) *SendmailEmailSender {
	sender := &SendmailEmailSender{
		cmd:    cmd,
		logger: log.NewNullLogger(),
		factory: execCmdSendmailCommandFactory,
	}
	if logger != nil {
		sender.logger = logger
	}
	return sender
}

func (s *SendmailEmailSender) Name() string {
	return "sendmail"
}

func (s *SendmailEmailSender) Send(from string, to string, subject string, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)
	cmd := s.factory(s.cmd, "-t")
	stdin, err := cmd.StdinPipe() // The pipe will be closed automatically after Wait sees the command exit.
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}
	_, err = m.WriteTo(stdin)            // with gomail we dont know how many bytes we have to write to check correctness
	if err = stdin.Close(); err != nil { // close directly to say sendmail to send message
		return errors.Wrap(err, "EmailSender[sendmail]: closing stdin")
	}
	if err = cmd.Wait(); err != nil {
		return errors.Wrap(err, "EmailSender[sendmail]: waiting")
	}
	return nil
}

func (s *SendmailEmailSender) Health() error {
	return nil
}

type SendmailCommandFactory func(name string, arg ...string) SendmailCommand

type SendmailCommand interface {
	Start() error
	StdinPipe() (io.WriteCloser, error)
	Wait() error
}

var execCmdSendmailCommandFactory SendmailCommandFactory = func(name string, arg ...string) SendmailCommand {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return &execCmdSendmailCommand{
		cmd: cmd,
	}
}

type execCmdSendmailCommand struct {
	cmd *exec.Cmd
}

func (s *execCmdSendmailCommand) Start() error {
	return s.cmd.Start()
}

func (s *execCmdSendmailCommand) StdinPipe() (io.WriteCloser, error) {
	return s.cmd.StdinPipe()
}

func (s *execCmdSendmailCommand) Wait() error {
	return s.cmd.Wait()
}
