package emailer

type EmailSender interface {
	Send(from string, to string, subject string, message string) error
}
