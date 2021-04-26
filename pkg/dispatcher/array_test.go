package dispatcher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayDriverNoError(t *testing.T) {
	d := NewArrayDriver()
	assert.NoError(t, d.Run())
	assert.NoError(t, d.Enqueue(func() error {
		return nil
	}))
	assert.NoError(t, d.Stop())
	assert.Len(t, d.Storage, 1)
}
