package leave_feedback

import (
	"github.com/mingalevme/feedbacker/internal/app"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/internal/app/repository"
	"github.com/mingalevme/feedbacker/pkg/errutils"
	"github.com/pkg/errors"
	"sync"
)

var ErrUnprocessableEntity = errors.New(repository.ErrUnprocessableEntity.Error())

type LeaveFeedback struct {
	env app.Env
	// wg is just for testing
	wg  sync.WaitGroup
}

func New(env app.Env) *LeaveFeedback {
	return &LeaveFeedback{
		env: env,
	}
}

func (s *LeaveFeedback) LeaveFeedback(input LeaveFeedbackData) (model.Feedback, error) {
	if err := input.Validate(); err != nil {
		return model.Feedback{}, err
	}
	data := convertLeaveFeedbackDataToAddFeedbackData(input)
	f, err := s.env.FeedbackRepository().Add(data)
	if errors.Is(err, repository.ErrUnprocessableEntity) {
		return f, ErrUnprocessableEntity
	}
	if err != nil {
		return f, err
	}
	s.wg.Add(1)
	go func() {
		defer func() {
			s.wg.Done()
			if r := recover(); r != nil {
				s.env.Logger().WithError(errutils.PanicToError(r)).Fatal("notifying")
			}
		}()
		if err := s.env.Notifier().Notify(f); err != nil {
			s.env.Logger().WithError(err).Error("notifying")
		}
	}()
	return f, nil
}

func convertLeaveFeedbackDataToAddFeedbackData(input LeaveFeedbackData) repository.AddFeedbackData {
	f := model.Feedback{
		Service: input.App,
		Edition: input.Edition,
		Text:    input.Body,
		Context: &model.Context{
			AppVersion:  input.AppVersion,
			AppBuild:    input.AppBuildNumber,
			OsName:      input.OsName,
			OsVersion:   input.OsVersion,
			DeviceBrand: input.Brand,
			DeviceModel: input.Model,
		},
		Customer: &model.Customer{
			Email:          input.Email,
			InstallationID: input.InstallationID,
		},
	}
	return repository.AddFeedbackData{Feedback: f}
}
