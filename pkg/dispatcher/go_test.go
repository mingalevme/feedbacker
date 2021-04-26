package dispatcher

import (
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestGoDriverNoError(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewGoDriver(l,1)
	assert.NoError(t, d.Run())
	var foo int
	assert.NoError(t, d.Enqueue(func() error {
		foo = 1
		return nil
	}))
	assert.Equal(t, 1, d.RunningProcessCount())
	assert.NoError(t, d.Stop())
	assert.Equal(t, 1, foo)
	assert.Equal(t, 0, d.RunningProcessCount())
	assert.Len(t, l.Storage(), 0)
}

func TestGoDriverError(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewGoDriver(l, 1)
	assert.NoError(t, d.Run())
	assert.NoError(t, d.Enqueue(func() error {
		return errors.New("TEST")
	}))
	assert.NoError(t, d.Stop())
	assert.Len(t, l.Storage(), 1)
}

func TestGoDriverPanicCatching(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewGoDriver(l, 2)
	assert.NoError(t, d.Run())
	assert.Equal(t, 0, d.RunningProcessCount())
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
	assert.Equal(t, 0, d.RunningProcessCount())
	assert.NoError(t, d.Stop())
	assert.Len(t, l.Storage(), 2)
}

func TestGoDriverOverflow(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewGoDriver(l, 1)
	assert.NoError(t, d.Run())
	assert.NoError(t, d.Enqueue(func() error {
		return nil
	}))
	assert.Error(t, MaxQueueSizeReached, d.Enqueue(func() error {
		return nil
	}))
	assert.NoError(t, d.Stop())
	assert.Len(t, l.Storage(), 0)
}

func TestGoDriverStopping(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewGoDriver(l, 1)
	assert.NoError(t, d.Run())
	assert.NoError(t, d.Stop())
	assert.Error(t, TaskQueueIsStopped, d.Enqueue(func() error {
		return nil
	}))
}

func TestGoDriverHealthNoError(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewGoDriver(l, 1)
	assert.NoError(t, d.Run())
	assert.NoError(t, d.Health())
	assert.NoError(t, d.Stop())
}

func TestGoHealthStopped(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewGoDriver(l, 1)
	assert.NoError(t, d.Run())
	assert.NoError(t, d.Stop())
	assert.Error(t, TaskQueueIsStopped, d.Health())
}

func TestGoHealthProcessesHaveReachedLimit(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewGoDriver(l, 1)
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
	assert.Contains(t, err.Error(), "process count has reached limit")
}

func TestGoHealthQueueIsCloseToLimit(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewGoDriver(l, 5)
	assert.NoError(t, d.Run())
	wg := sync.WaitGroup{}
	wg.Add(1)
	for i := 1; i <= 4; i++ {
		assert.NoError(t, d.Enqueue(func() error {
			wg.Wait()
			return nil
		}))
	}
	assert.Equal(t, 4, d.RunningProcessCount())
	err := d.Health()
	wg.Done()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "running process count has reached 80.00% of capacity")
	assert.NoError(t, d.Stop())
}
