package dispatcher

import (
	"fmt"
	"github.com/mingalevme/feedbacker/pkg/errutils"
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
)

type GoDriver struct {
	logger                 log.Logger
	wg                     sync.WaitGroup
	maxRunningProcessCount int
	runningProcessCount    uint64 // atomic!
	enqueueingMutex        sync.Mutex
	isRunning              bool
}

func NewGoDriver(logger log.Logger, maxRunningProcessCount int) *GoDriver {
	if maxRunningProcessCount < 1 {
		panic("maxRunningProcessCount must be positive")
	}
	return &GoDriver{
		logger:                 logger,
		maxRunningProcessCount: maxRunningProcessCount,
		isRunning:              true,
	}
}

func (s *GoDriver) Name() string {
	return "go"
}

func (s *GoDriver) Enqueue(t Task) error {
	if !s.isRunning {
		return TaskQueueIsStopped
	}
	s.enqueueingMutex.Lock()
	if int(atomic.LoadUint64(&s.runningProcessCount)) > s.maxRunningProcessCount {
		return MaxQueueSizeReached
	}
	s.wg.Add(1)
	s.incRunningProcessCount()
	s.enqueueingMutex.Unlock()
	go func() {
		// panic catching
		defer func() {
			defer s.wg.Done()
			s.decRunningProcessCount()
			if r := recover(); r != nil {
				s.logger.WithError(errutils.PanicToError(r)).Fatal("dispatcher[go]: Panic while processing a task")
			}
		}()
		if err := t(); err != nil {
			s.logger.WithError(err).Error("Error while processing a task")
		}
	}()
	return nil
}

func (s *GoDriver) incRunningProcessCount() {
	atomic.AddUint64(&s.runningProcessCount, 1)
}

func (s *GoDriver) decRunningProcessCount() {
	atomic.AddUint64(&s.runningProcessCount, ^uint64(0))
}

func (s *GoDriver) RunningProcessCount() int {
	return int(atomic.LoadUint64(&s.runningProcessCount))
}

func (s *GoDriver) Run() error {
	return nil
}

func (s *GoDriver) Stop() error {
	s.isRunning = false
	fmt.Println("Dispatcher[go]/worker: Stopping workers gracefully")
	s.wg.Wait()
	fmt.Println("Dispatcher[go]/worker: Workers have been stopped")
	return nil
}

func (s *GoDriver) Health() error {
	if !s.isRunning {
		return TaskQueueIsStopped
	}
	runningProcessCount := s.RunningProcessCount()
	if runningProcessCount == s.maxRunningProcessCount {
		return errors.Errorf("dispatcher[go]: running process count has reached limit (%d)", s.maxRunningProcessCount)
	}
	if runningProcessCount > s.maxRunningProcessCount*3/4 {
		return errors.Errorf("dispatcher[go]: running process count has reached %0.2f%% of capacity", 100*float64(runningProcessCount)/float64(s.maxRunningProcessCount))
	}
	return nil
}
