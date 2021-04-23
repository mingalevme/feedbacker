package dispatcher

import "github.com/mingalevme/feedbacker/pkg/log"

type NullDriver struct {
	logger log.Logger
}

func NewNullDriver(logger log.Logger) *NullDriver {
	return &NullDriver{
		logger: logger,
	}
}

func (s *NullDriver) Name() string {
	return "null"
}

func (s *NullDriver) Enqueue(t Task) error {
	return nil
}

func (s *NullDriver) Run() error {
	return nil
}

func (s *NullDriver) Stop() error {
	return nil
}

func (s *NullDriver) Health() error {
	return nil
}
