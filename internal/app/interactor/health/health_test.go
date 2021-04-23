package health

import (
	"github.com/mingalevme/feedbacker/internal/app"
	"github.com/mingalevme/feedbacker/pkg/timeutils"
	"github.com/mingalevme/feedbacker/test"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHealthSuccess(t *testing.T) {
	now := time.Now()
	timeutils.SetTestNow(now)
	env := testEnv()
	i := New(env)
	h := i.Health()
	assert.Equal(t, HealthStatusPass, h.Status)
	assert.Equal(t, "", h.Output)
	assert.Len(t, h.Details, 3)
	//
	repoCompName := "repository/array"
	assert.Len(t, h.Details[repoCompName], 1)
	assert.Equal(t, HealthStatusPass, h.Details[repoCompName][0].Status)
	assert.Equal(t, now.UTC().Format(time.RFC3339), h.Details[repoCompName][0].Time)
	assert.Equal(t, "", h.Details[repoCompName][0].Output)
	//
	notifierCompName := "notifier/array"
	assert.Len(t, h.Details[notifierCompName], 1)
	assert.Equal(t, HealthStatusPass, h.Details[notifierCompName][0].Status)
	assert.Equal(t, now.UTC().Format(time.RFC3339), h.Details[notifierCompName][0].Time)
	assert.Equal(t, "", h.Details[notifierCompName][0].Output)
	//
	dispatcherCompName := "dispatcher/sync"
	assert.Len(t, h.Details[dispatcherCompName], 1)
	assert.Equal(t, HealthStatusPass, h.Details[dispatcherCompName][0].Status)
	assert.Equal(t, now.UTC().Format(time.RFC3339), h.Details[dispatcherCompName][0].Time)
	assert.Equal(t, "", h.Details[dispatcherCompName][0].Output)
}

func TestHealthRepoError(t *testing.T) {
	env := testEnv()
	r := &test.HealthErrorRepository{
		Err: errors.New("TEST"),
	}
	env.SetFeedbackRepository(r)
	i := New(env)
	h := i.Health()
	//
	repoCompName := "repository/error"
	//
	assert.Equal(t, HealthStatusFail, h.Status)
	assert.Equal(t, repoCompName+": TEST", h.Output)
	//
	assert.Len(t, h.Details[repoCompName], 1)
	assert.Equal(t, HealthStatusFail, h.Details[repoCompName][0].Status)
	assert.Equal(t, "TEST", h.Details[repoCompName][0].Output)
}

func TestHealthNotifierError(t *testing.T) {
	env := testEnv()
	n := &test.HealthErrorNotifier{
		Err: errors.New("TEST"),
	}
	env.SetNotifier(n)
	i := New(env)
	h := i.Health()
	//
	compName := "notifier/error"
	//
	assert.Equal(t, HealthStatusFail, h.Status)
	assert.Equal(t, compName + ": TEST", h.Output)
	//
	assert.Len(t, h.Details[compName], 1)
	assert.Equal(t, HealthStatusFail, h.Details[compName][0].Status)
	assert.Equal(t, "TEST", h.Details[compName][0].Output)
}

func TestHealthDispatcherError(t *testing.T) {
	env := testEnv()
	d := &test.HealthErrorDispatcher{
		Err: errors.New("TEST"),
	}
	env.SetDispatcher(d)
	i := New(env)
	h := i.Health()
	//
	compName := "dispatcher/error"
	//
	assert.Equal(t, HealthStatusFail, h.Status)
	assert.Equal(t, compName+": TEST", h.Output)
	//
	assert.Len(t, h.Details[compName], 1)
	assert.Equal(t, HealthStatusFail, h.Details[compName][0].Status)
	assert.Equal(t, "TEST", h.Details[compName][0].Output)
}

func testEnv(values ...map[string]string) *app.Container {
	base := map[string]string{
		"PERSISTENCE_DRIVER": "array",
		"NOTIFIER_CHANNEL":   "array",
		"DISPATCHER_DRIVER":  "sync",
	}
	for _, m := range values {
		for k, v := range m {
			base[k] = v
		}
	}
	return test.NewEnv(base)
}
