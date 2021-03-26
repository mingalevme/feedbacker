package services

import (
	"github.com/mingalevme/feedbacker/domain/feedback"
)

type LeaveFeedbackData struct {
	Service  string                    `json:"service,omitempty"`
	Edition  string                    `json:"edition,omitempty"`
	Text     string                    `json:"text,omitempty"`
	Customer LeaveFeedbackCustomerData `json:"customer,omitempty"`
	Context  LeaveFeedbackContextData  `json:"context,omitempty"`
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
}

func NewLeaveFeedbackService(repository feedback.Repository) LeaveFeedbackService {
	service := leaveFeedbackService{
		repository: repository,
	}
	return service
}

func (s leaveFeedbackService) Handle(data LeaveFeedbackData) (feedback.FeedbackData, error) {
	serviceValue, err := feedback.NewServiceValue(data.Service)
	if err != nil {
		return nil, err
	}
	editionValue, err := feedback.NewEditionValue(data.Edition)
	if err != nil {
		return nil, err
	}
	textValue, err := feedback.NewTextValue(data.Text)
	if err != nil {
		return nil, err
	}
	customer, err := feedback.NewCustomer(data.Customer.Email, data.Customer.InstallationID)
	if err != nil {
		return nil, err
	}
	context, err := feedback.NewContext(data.Context.AppVersion, data.Context.AppBuild, data.Context.OsName, data.Context.OsVersion, data.Context.DeviceBrand, data.Context.DeviceModel)
	if err != nil {
		return nil, err
	}
	f, err := s.repository.Add(serviceValue, editionValue, textValue, customer, context)
	if err != nil {
		return nil, err
	}
	return f, nil
}
