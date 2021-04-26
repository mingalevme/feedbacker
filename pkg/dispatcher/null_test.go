package dispatcher

import (
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNullDriverNoError(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewNullDriver(l)
	assert.NoError(t, d.Run())
	var foo int
	assert.NoError(t, d.Enqueue(func() error {
		foo = 1
		return nil
	}))
	assert.NoError(t, d.Stop())
	assert.Equal(t, 0, foo)
	assert.Len(t, l.Storage(), 0)
}

func TestNullDriverError(t *testing.T) {
	l := log.NewArrayLogger(log.LevelError)
	d := NewNullDriver(l)
	assert.NoError(t, d.Run())
	assert.NoError(t, d.Enqueue(func() error {
		return errors.New("TEST")
	}))
	assert.NoError(t, d.Stop())
	assert.Len(t, l.Storage(), 0)
}
