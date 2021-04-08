package interactor

import (
	"github.com/mingalevme/feedbacker/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPing(t *testing.T) {
	i := New(test.NewEnv(nil))
	assert.Equal(t, "pong", i.Ping())
}
