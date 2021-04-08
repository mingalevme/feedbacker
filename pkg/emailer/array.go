package emailer

type ArrayEmailItem struct {
	From    string
	To      string
	Subject string
	Message string
}

type ArrayEmailSender struct {
	Storage []ArrayEmailItem
}

func NewArrayEmailSender() *ArrayEmailSender {
	sender := &ArrayEmailSender{
		Storage: []ArrayEmailItem{},
	}
	return sender
}

func (s *ArrayEmailSender) Send(from string, to string, subject string, message string) error {
	s.Storage = append(s.Storage, ArrayEmailItem{
		From:    from,
		To:      to,
		Subject: subject,
		Message: message,
	})
	return nil
}
