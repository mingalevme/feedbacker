package dispatcher

import (
	"fmt"
	"github.com/mingalevme/feedbacker/pkg/errutils"
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
)

type ChanDriver struct {
	queue              chan Task
	queueSize          uint64
	queueMaxSize       int
	isRunning          bool
	workerCount        int
	quit               chan bool
	logger             log.Logger
	workerWG           sync.WaitGroup
	runningWorkerCount uint64
}

func NewChanDriver(logger log.Logger, queueMaxSize int, workersCount int) *ChanDriver {
	if workersCount < 1 {
		panic(errors.New("workers count is invalid"))
	}
	return &ChanDriver{
		queue:        make(chan Task, queueMaxSize),
		queueMaxSize: queueMaxSize,
		isRunning:    true,
		workerCount:  workersCount,
		quit:         make(chan bool),
		logger:       logger,
	}
}

func (s *ChanDriver) Name() string {
	return "chan"
}

func (s *ChanDriver) Enqueue(t Task) error {
	if !s.isRunning {
		return TaskQueueIsStopped
	}
	// Do not block if limit has been reached
	select {
	case s.queue <- t:
		atomic.AddUint64(&s.queueSize, 1)
		return nil
	default:
		return MaxQueueSizeReached
	}
}

func (s *ChanDriver) Run() error {
	wg := sync.WaitGroup{}
	for i := 1; i < s.workerCount+1; i++ {
		s.workerWG.Add(1)
		wg.Add(1)
		fmt.Printf("Dispatcher[chan]: Starting worker %d/%d\n", i, s.workerCount)
		go func() {
			wg.Done()
			s.work()
		}()
	}
	wg.Wait()
	return nil
}

func (s *ChanDriver) work() {
	defer s.workerWG.Done()
	s.incRunningWorkerCount()
	defer s.decRunningWorkerCount()
	for {
		select {
		case <-s.quit:
			fmt.Println("Dispatcher[chan]/worker: Halt signal received, stopping the worker")
			return
		case t, ok := <-s.queue:
			if !ok { // channel is closed
				s.logger.Debug("worker: channel is closed, exiting")
				fmt.Println("Dispatcher[chan]/worker: Stop signal received, stopping the worker")
				return
			}
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				// panic catching
				defer func() {
					atomic.AddUint64(&s.queueSize, ^uint64(0))
					defer wg.Done()
					if r := recover(); r != nil {
						s.logger.WithError(errutils.PanicToError(r)).Fatal("dispatcher[chan]: Error while processing the task")
					}
				}()
				if err := t(); err != nil {
					s.logger.WithError(err).Error("Error while processing a task")
				}
			}()
			wg.Wait()
		}
	}
}

func (s *ChanDriver) Stop() error {
	fmt.Println("Dispatcher[chan]: Stopping workers gracefully ...")
	s.isRunning = false
	close(s.queue)
	s.workerWG.Wait()
	fmt.Println("Dispatcher[chan]: Workers have been stopped")
	return nil
}

func (s *ChanDriver) Halt() error {
	fmt.Printf("Dispatcher[chan]: Halting workers")
	s.quit <- false
	s.workerWG.Wait()
	fmt.Printf("Dispatcher[chan]: Workers have been halted")
	return nil
}

func (s *ChanDriver) incRunningWorkerCount() {
	atomic.AddUint64(&s.runningWorkerCount, 1)
}

func (s *ChanDriver) decRunningWorkerCount() {
	atomic.AddUint64(&s.runningWorkerCount, ^uint64(0))
}

func (s *ChanDriver) RunningWorkerCount() int {
	return int(atomic.LoadUint64(&s.runningWorkerCount))
}

func (s *ChanDriver) QueueSize() int {
	return int(atomic.LoadUint64(&s.queueSize))
}

func (s *ChanDriver) Health() error {
	if !s.isRunning {
		return TaskQueueIsStopped
	}
	if s.workerCount != int(s.runningWorkerCount) {
		return errors.Errorf("dispatcher[chan]: the number of running process (%d) is different than expected (%d)", s.runningWorkerCount, s.workerCount)
	}
	queueSize := s.QueueSize()
	if queueSize == s.queueMaxSize {
		return errors.Errorf("dispatcher[chan]: queue has reached limit (%d)", s.queueMaxSize)
	}
	if queueSize > s.queueMaxSize*3/4 {
		return errors.Errorf("dispatcher[chan]: queue has reached %0.2f%% of capacity", 100*float64(queueSize)/float64(s.queueMaxSize))
	}
	return nil
}
