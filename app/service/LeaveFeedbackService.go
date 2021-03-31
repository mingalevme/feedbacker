package service

import (
	"github.com/mingalevme/feedbacker/domain/feedback"
	"github.com/mingalevme/feedbacker/infrastructure/log"
	"github.com/pkg/errors"
)

type LeaveFeedbackData struct {
	Service  string                     `json:"service,omitempty"`
	Edition  string                     `json:"edition,omitempty"`
	Text     string                     `json:"text,omitempty"`
	Customer *LeaveFeedbackCustomerData `json:"customer,omitempty"`
	Context  *LeaveFeedbackContextData  `json:"context,omitempty"`
}

type LeaveFeedbackCustomerData struct {
	Email          *string `json:"email,omitempty"`
	InstallationID *string `json:"installationId,omitempty"`
}

type LeaveFeedbackContextData struct {
	AppVersion  *string `json:"appVersion,omitempty"`
	AppBuild    *string `json:"appBuild,omitempty"`
	OsName      *string `json:"osName,omitempty"`
	OsVersion   *string `json:"osVersion,omitempty"`
	DeviceBrand *string `json:"deviceBrand,omitempty"`
	DeviceModel *string `json:"deviceModel,omitempty"`
}

type LeaveFeedbackService interface {
	Handle(data LeaveFeedbackData) (feedback.FeedbackData, error)
}

type leaveFeedbackService struct {
	repository feedback.Repository
	logger     log.Logger
}

func NewLeaveFeedbackService(repository feedback.Repository, logger log.Logger) LeaveFeedbackService {
	service := leaveFeedbackService{
		repository: repository,
		logger:     logger,
	}
	return &service
}

func (s *leaveFeedbackService) Handle(data LeaveFeedbackData) (feedback.FeedbackData, error) {
	serviceValue, err := feedback.NewServiceValue(data.Service)
	if err != nil {
		return nil, errors.Wrap(ErrUnprocessableEntity, err.Error())
	}
	editionValue, err := feedback.NewEditionValue(data.Edition)
	if err != nil {
		return nil, errors.Wrap(ErrUnprocessableEntity, err.Error())
	}
	textValue, err := feedback.NewTextValue(data.Text)
	if err != nil {
		return nil, errors.Wrap(ErrUnprocessableEntity, err.Error())
	}
	var customer feedback.Customer
	if data.Customer != nil {
		customer, err = feedback.NewCustomer(data.Customer.Email, data.Customer.InstallationID)
		if err != nil {
			return nil, errors.Wrap(ErrUnprocessableEntity, err.Error())
		}
	}
	var context feedback.Context
	if data.Context != nil {
		context, err = feedback.NewContext(data.Context.AppVersion, data.Context.AppBuild, data.Context.OsName, data.Context.OsVersion, data.Context.DeviceBrand, data.Context.DeviceModel)
		if err != nil {
			return nil, errors.Wrap(ErrUnprocessableEntity, err.Error())
		}
	}
	f, err := s.repository.Add(serviceValue, editionValue, textValue, customer, context)
	if err != nil {
		s.logger.WithError(err).WithField("data", data).Error("Error while adding feedback")
		return nil, err
	}
	return f, nil
}
