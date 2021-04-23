// +build testing

package notifier

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogLogging(t *testing.T) {
	logger := log.NewArrayLogger()
	n := NewLogNotifier(logger, log.LevelDebug)
	err := n.Notify(model.MakeFeedback())
	assert.NoError(t, err)
	assert.Len(t, logger.Storage, 1)
}
