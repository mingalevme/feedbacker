package dispatcher

import (
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSyncDriverNoError(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewSyncDriver(l)
	assert.NoError(t, d.Run())
	var foo int
	assert.NoError(t, d.Enqueue(func() error {
		foo = 1
		return nil
	}))
	assert.NoError(t, d.Stop())
	assert.Equal(t, 1, foo)
	assert.Len(t, l.Storage(), 0)
}

func TestSyncDriverError(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewSyncDriver(l)
	assert.NoError(t, d.Run())
	assert.NoError(t, d.Enqueue(func() error {
		return errors.New("TEST")
	}))
	assert.NoError(t, d.Stop())
	assert.Len(t, l.Storage(), 1)
}
