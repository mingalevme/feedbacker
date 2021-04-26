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
