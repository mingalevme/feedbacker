package emailer

type EmailSender interface {
	Name() string
	Health() error
	Send(from string, to string, subject string, message string) error
}
