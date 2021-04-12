package emailer

type NullEmailSender struct{}

func (s *NullEmailSender) Name() string {
	return "null"
}

func (s *NullEmailSender) Health() error {
	return nil
}

func NewNullEmailSender() *NullEmailSender {
	return &NullEmailSender{}
}

func (s *NullEmailSender) Send(from string, to string, subject string, message string) error {
	return nil
}
