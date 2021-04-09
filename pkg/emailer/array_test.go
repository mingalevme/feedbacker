package emailer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayName(t *testing.T) {
	arrayEmailSender := NewArrayEmailSender()
	assert.Equal(t, "array", arrayEmailSender.Name())
}

func TestArrayHealthNoError(t *testing.T) {
	arrayEmailSender := NewArrayEmailSender()
	err := arrayEmailSender.Health()
	assert.NoError(t, err)
}
