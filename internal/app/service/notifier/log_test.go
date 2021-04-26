// +build testing

package notifier

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogLogging(t *testing.T) {
	l := log.NewArrayLogger(log.LevelDebug)
	n := NewLogNotifier(l, log.LevelDebug)
	err := n.Notify(model.MakeFeedback())
	assert.NoError(t, err)
	assert.Len(t, l.Storage(), 1)
}
