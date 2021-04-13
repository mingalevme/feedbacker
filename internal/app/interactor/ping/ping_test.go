package ping

import (
	"github.com/mingalevme/feedbacker/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPing(t *testing.T) {
	ping := New(test.NewEnv(nil))
	assert.Equal(t, "pong", ping.Ping())
}
