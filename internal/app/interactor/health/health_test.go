package health

import (
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
	env := test.NewEnv(map[string]string{
		"PERSISTENCE_DRIVER": "array",
		"NOTIFIER_CHANNEL":   "array",
	})
	i := New(env)
	h := i.Health()
	assert.Equal(t, HealthStatusPass, h.Status)
	assert.Equal(t, "", h.Output)
	assert.Len(t, h.Details, 2)
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
}

func TestHealthRepoError(t *testing.T) {
	env := test.NewEnv(map[string]string{
		"PERSISTENCE_DRIVER": "array",
		"NOTIFIER_CHANNEL":   "array",
	})
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
	assert.Equal(t, repoCompName + ": TEST", h.Output)
	//
	assert.Len(t, h.Details[repoCompName], 1)
	assert.Equal(t, HealthStatusFail, h.Details[repoCompName][0].Status)
	assert.Equal(t, "TEST", h.Details[repoCompName][0].Output)
}

func TestHealthNotifierError(t *testing.T) {
	env := test.NewEnv(map[string]string{
		"PERSISTENCE_DRIVER": "array",
		"NOTIFIER_CHANNEL":   "array",
	})
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
