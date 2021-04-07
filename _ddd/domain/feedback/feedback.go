// https://google.github.io/styleguide/jsoncstyleguide.xml

package feedback

import "time"

type Feedback interface {
	GetID() FeedbackID
	GetService() ServiceValue
	GetEdition() EditionValue
	GetText() TextValue
	GetCustomer() Customer
	GetContext() Context
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetService(value string)
	SetEdition(value string)
	SetText(value string)
	SetCustomer(customer Customer)
	SetContext(context Context)
	Equals(value interface{}) bool
	MarshalJSON() ([]byte, error)
}

type feedback FeedbackData

//type feedback struct {
//	FeedbackData
//	//id        feedbackId   `json:"id,omitempty"`
//	//service   serviceValue `json:"service,omitempty"`
//	//edition   editionValue `json:"edition,omitempty"`
//	//text      textValue    `json:"text,omitempty"`
//	//customer  Customer     `json:"customer,omitempty"`
//	//context   context      `json:"context,omitempty"`
//	//createdAt time.Time    `json:"created_at,omitempty"`
//	//updatedAt time.Time    `json:"updated_at,omitempty"`
//}

func (f *feedback) GetID() FeedbackID {
	return &feedbackID{
		value: f.Id,
	}
}

func (f *feedback) GetService() ServiceValue {
	return &serviceValue{
		abstractStringValueObject{
			value: f.Service,
		},
	}
}

func (f *feedback) GetEdition() EditionValue {
	return &editionValue{
		abstractStringValueObject{
			value: f.Edition,
		},
	}
}

func (f *feedback) GetText() TextValue {
	return &textValue{
		abstractStringValueObject{
			value: f.Text,
		},
	}
}

func (f *feedback) GetCustomer() Customer {
	return &customer{
		f.Customer,
	}
}

func (f *feedback) GetContext() *Context {
	return f.Context
}

func (f *feedback) GetCreatedAt() time.Time {
	return f.CreatedAt
}

func (f *feedback) GetUpdatedAt() time.Time {
	return f.UpdatedAt
}

func (f *feedback) Equals(value interface{}) bool {
	target, ok := value.(int)
	return ok && f.Id == target
}

func (f *feedback) SetService(s string) {
	f.Service = s
}

func (f *feedback) SetEdition(edition string) {
	f.Edition = edition
}

func (f *feedback) SetText(text string) {
	f.Text = text
}

func (f *feedback) SetCustomer(customer *Customer) {
	f.Customer = customer
}

func (f *feedback) SetContext(context *Context) {
	f.Context = context
}

//func NewFeedback(service serviceValue, edition editionValue, text textValue, context context) *feedback {
//	return &feedback{
//		service: service,
//		edition: edition,
//		text:    text,
//		context: context,
//	}
//}
