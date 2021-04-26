package dispatcher

import (
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestChanDriverNoError(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewChanDriver(l, 1, 1)
	assert.NoError(t, d.Run())
	var foo int
	assert.NoError(t, d.Enqueue(func() error {
		foo = 1
		return nil
	}))
	assert.NoError(t, d.Stop())
	assert.Equal(t, 1, foo)
	assert.Equal(t, uint64(0), d.queueSize)
	assert.Len(t, l.Storage(), 0)
}

func TestChanDriverError(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewChanDriver(l, 1, 1)
	assert.NoError(t, d.Run())
	assert.NoError(t, d.Enqueue(func() error {
		return errors.New("TEST")
	}))
	assert.NoError(t, d.Stop())
	assert.Equal(t, uint64(0), d.queueSize)
	assert.Len(t, l.Storage(), 1)
}

func TestChanDriverWorkerCount(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewChanDriver(l, 1, 2)
	assert.Equal(t, 0, d.RunningWorkerCount())
	assert.NoError(t, d.Run())
	assert.Equal(t, 2, d.RunningWorkerCount())
	assert.NoError(t, d.Stop())
	assert.Equal(t, 0, d.RunningWorkerCount())
}

func TestChanDriverPanicCatching(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewChanDriver(l, 1, 1)
	assert.NoError(t, d.Run())
	assert.Equal(t, 1, d.RunningWorkerCount())
	counter := 0
	wg := sync.WaitGroup{}
	wg.Add(2)
	assert.NoError(t, d.Enqueue(func() error {
		defer wg.Done()
		counter++
		panic("TEST")
	}))
	assert.NoError(t, d.Enqueue(func() error {
		defer wg.Done()
		counter++
		panic("TEST")
	}))
	wg.Wait()
	assert.Equal(t, 2, counter)
	assert.Equal(t, 1, d.RunningWorkerCount())
	assert.NoError(t, d.Stop())
	assert.Equal(t, uint64(0), d.queueSize)
	assert.Len(t, l.Storage(), 2)
}

func TestChanDriverOverflow(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewChanDriver(l, 1, 1)
	assert.NoError(t, d.Enqueue(func() error {
		return nil
	}))
	assert.Error(t, MaxQueueSizeReached, d.Enqueue(func() error {
		return nil
	}))
	assert.NoError(t, d.Stop())
}

func TestChanDriverStopping(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewChanDriver(l, 1, 1)
	assert.NoError(t, d.Run())
	assert.NoError(t, d.Stop())
	assert.Error(t, TaskQueueIsStopped, d.Enqueue(func() error {
		return nil
	}))
}

func TestChanHealthNoError(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewChanDriver(l, 1, 1)
	assert.NoError(t, d.Run())
	assert.NoError(t, d.Health())
	assert.NoError(t, d.Stop())
}

func TestChanHealthStopped(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewChanDriver(l, 1, 1)
	assert.NoError(t, d.Run())
	assert.NoError(t, d.Stop())
	assert.Error(t, TaskQueueIsStopped, d.Health())
}

func TestChanHealthQueueHasReachedLimit(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewChanDriver(l, 1, 1)
	assert.NoError(t, d.Run())
	wg := sync.WaitGroup{}
	wg.Add(1)
	assert.NoError(t, d.Enqueue(func() error {
		wg.Wait()
		return nil
	}))
	err := d.Health()
	wg.Done()
	assert.NoError(t, d.Stop())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "queue has reached limit")
}

func TestChanHealthQueueIsCloseToLimit(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewChanDriver(l, 5, 1)
	assert.NoError(t, d.Run())
	wg := sync.WaitGroup{}
	wg.Add(1)
	for i := 1; i <= 4; i++ {
		assert.NoError(t, d.Enqueue(func() error {
			wg.Wait()
			return nil
		}))
	}
	assert.Equal(t, 4, d.QueueSize())
	err := d.Health()
	wg.Done()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "queue has reached 80.00% of capacity")
	assert.NoError(t, d.Stop())
}
