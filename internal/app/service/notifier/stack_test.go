package notifier

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStackSendNoError(t *testing.T) {
	arr1 := NewArrayNotifier(nil)
	arr2 := NewArrayNotifier(nil)
	stack := NewStackNotifier()
	stack.Add(arr1)
	stack.Add(arr2)
	_ = stack.Notify(model.MakeFeedback())
	assert.Len(t, arr1.Storage, 1)
	assert.Len(t, arr2.Storage, 1)
}

func TestStackSendErr(t *testing.T) {
	n1 := NewErrNotifier(errors.New("Err1"))
	n2 := NewErrNotifier(errors.New("Err2"))
	stack := NewStackNotifier()
	stack.Add(n1)
	stack.Add(n2)
	err := stack.Notify(model.MakeFeedback())
	assert.Error(t, err)
	assert.Equal(t, "error: Err1; error: Err2", err.Error())
}
