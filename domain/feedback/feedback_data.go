package feedback

import "time"

type FeedbackData interface {
	GetID() FeedbackId
	GetService() ServiceValue
	GetEdition() EditionValue
	GetText() TextValue
	GetCustomer() Customer
	GetContext() Context
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	Equals(value interface{}) bool
}

type feedbackData struct {
	Id        FeedbackId  	`json:"id,omitempty"`
	Service   ServiceValue	`json:"service,omitempty"`
	Edition   EditionValue	`json:"edition,omitempty"`
	Text      TextValue		`json:"text,omitempty"`
	Customer  Customer		`json:"customer,omitempty"`
	Context   Context		`json:"context,omitempty"`
	CreatedAt time.Time		`json:"created_at,omitempty"`
	UpdatedAt time.Time		`json:"updated_at,omitempty"`
}

func (f feedbackData) GetID() FeedbackId {
	return f.Id
}

func (f feedbackData) GetService() ServiceValue {
	return f.Service
}

func (f feedbackData) GetEdition() EditionValue {
	return f.Edition
}

func (f feedbackData) GetText() TextValue {
	return f.Text
}

func (f feedbackData) GetCustomer() Customer {
	return f.Customer
}

func (f feedbackData) GetContext() Context {
	return f.Context
}

func (f feedbackData) GetCreatedAt() time.Time {
	return f.CreatedAt
}

func (f feedbackData) GetUpdatedAt() time.Time {
	return f.UpdatedAt
}

func (f feedbackData) Equals(value interface{}) bool {
	target, ok := value.(FeedbackData)
	return ok && f.Id.Equals(target.GetID())
}

func NewFeedbackData(f Feedback) FeedbackData {
	return feedbackData{
		Id:        f.GetID(),
		Service:   f.GetService(),
		Edition:   f.GetEdition(),
		Text:      f.GetText(),
		Customer:  f.GetCustomer(),
		Context:   f.GetContext(),
		CreatedAt: f.GetCreatedAt(),
		UpdatedAt: f.GetUpdatedAt(),
	}
}

