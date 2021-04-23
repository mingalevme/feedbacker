package dispatcher

import (
	"github.com/mingalevme/feedbacker/pkg/log"
)

type SyncDriver struct {
	logger log.Logger
}

func NewSyncDriver(logger log.Logger) *SyncDriver {
	return &SyncDriver{
		logger: logger,
	}
}

func (s *SyncDriver) Name() string {
	return "sync"
}

func (s *SyncDriver) Enqueue(t Task) error {
	if err := t(); err != nil {
		s.logger.WithError(err).Error("Error while processing a task")
	}
	return nil
}

func (s *SyncDriver) Run() error {
	return nil
}

func (s *SyncDriver) Stop() error {
	return nil
}

func (s *SyncDriver) Health() error {
	return nil
}
