package dispatcher

import (
	"github.com/mingalevme/feedbacker/pkg/log"
)

type GoDriver struct {
	logger log.Logger
}

func NewGoDriver(logger log.Logger) *GoDriver {
	return &GoDriver{
		logger: logger,
	}
}

func (s *GoDriver) Name() string {
	return "go"
}

func (s *GoDriver) Enqueue(t Task) error {
	go func() {
		_ = t()
	}()
	return nil
}

func (s *GoDriver) Run() error {
	return nil
}

func (s *GoDriver) Stop() error {
	return nil
}

func (s *GoDriver) Health() error {
	return nil
}
