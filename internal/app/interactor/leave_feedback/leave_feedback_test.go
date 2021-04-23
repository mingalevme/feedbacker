package leave_feedback

import (
	"github.com/mingalevme/feedbacker/internal/app"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/internal/app/repository"
	"github.com/mingalevme/feedbacker/internal/app/service/notifier"
	"github.com/mingalevme/feedbacker/pkg/dispatcher"
	"github.com/mingalevme/feedbacker/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLeaveFeedbackSuccess(t *testing.T) {
	env := testEnv()
	i := New(env)
	r, _ := env.FeedbackRepository().(*repository.ArrayFeedbackRepository)
	assert.Len(t, r.Storage, 0)
	n, _ := env.Notifier().(*notifier.ArrayNotifier)
	assert.Len(t, n.Storage, 0)
	d, _ := env.Dispatcher().(*dispatcher.ArrayDriver)
	assert.Len(t, d.Storage, 0)
	f1 := model.MakeFeedback()
	payload := LeaveFeedbackData{
		App:            f1.Service,
		AppVersion:     f1.Context.AppVersion,
		AppBuildNumber: f1.Context.AppBuild,
		Edition:        f1.Edition,
		Body:           f1.Text,
		Brand:          f1.Context.DeviceBrand,
		Model:          f1.Context.DeviceModel,
		OsName:         f1.Context.OsName,
		OsVersion:      f1.Context.OsVersion,
		Email:          f1.Customer.Email,
		InstallationID: f1.Customer.InstallationID,
	}
	f2, err := i.LeaveFeedback(payload)
	assert.NoError(t, err)
	assert.Len(t, r.Storage, 1)
	assert.Len(t, n.Storage, 0)
	assert.Len(t, d.Storage, 1)
	assert.Equal(t, f1.Service, f2.Service)
	assert.Equal(t, f1.Context.AppVersion, f2.Context.AppVersion)
	assert.Equal(t, f1.Context.AppBuild, f2.Context.AppBuild)
	assert.Equal(t, f1.Edition, f2.Edition)
	assert.Equal(t, f1.Text, f2.Text)
	assert.Equal(t, f1.Context.DeviceBrand, f2.Context.DeviceBrand)
	assert.Equal(t, f1.Context.DeviceModel, f2.Context.DeviceModel)
	assert.Equal(t, f1.Context.OsName, f2.Context.OsName)
	assert.Equal(t, f1.Context.OsVersion, f2.Context.OsVersion)
	assert.Equal(t, f1.Customer.Email, f2.Customer.Email)
	assert.Equal(t, f1.Customer.InstallationID, f2.Customer.InstallationID)
	_ = d.Storage[0]()
	assert.Len(t, n.Storage, 1)
}

func TestLeaveFeedbackUnprocessableEntity(t *testing.T) {
	env := testEnv()
	i := New(env)
	r, _ := env.FeedbackRepository().(*repository.ArrayFeedbackRepository)
	assert.Len(t, r.Storage, 0)
	d, _ := env.Dispatcher().(*dispatcher.ArrayDriver)
	assert.Len(t, d.Storage, 0)
	n, _ := env.Notifier().(*notifier.ArrayNotifier)
	assert.Len(t, n.Storage, 0)
	payload := LeaveFeedbackData{
		App: "",
	}
	_, err := i.LeaveFeedback(payload)
	assert.ErrorIs(t, err, ErrUnprocessableEntity)
	assert.Len(t, r.Storage, 0)
	assert.Len(t, d.Storage, 0)
	assert.Len(t, n.Storage, 0)
}

func testEnv(values ...map[string]string) *app.Container {
	base := map[string]string{
		"PERSISTENCE_DRIVER": "array",
		"NOTIFIER_CHANNEL":   "array",
		"DISPATCHER_DRIVER":  "array",
	}
	for _, m := range values {
		for k, v := range m {
			base[k] = v
		}
	}
	return test.NewEnv(base)
}
