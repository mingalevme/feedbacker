// https://google.github.io/styleguide/jsoncstyleguide.xml

package feedback

import (
	"github.com/pkg/errors"
)

var IdIsInvalid = errors.New("id is invalid")
var ServiceIsInvalid = errors.New("service is invalid")
var EditionIsInvalid = errors.New("edition is invalid")
var TextIsInvalid = errors.New("text is invalid")
var TextIsEmpty = errors.New("text is empty")

type Feedback interface {
	FeedbackData
	SetService(value ServiceValue)
	SetEdition(value EditionValue)
	SetText(value TextValue)
	SetCustomer(customer Customer)
	SetContext(context Context)
}

type feedback struct {
	feedbackData
	//id        feedbackId           `json:"id,omitempty"`
	//service   serviceValue `json:"service,omitempty"`
	//edition   editionValue `json:"edition,omitempty"`
	//text      textValue    `json:"text,omitempty"`
	//customer  Customer     `json:"customer,omitempty"`
	//context   context      `json:"context,omitempty"`
	//createdAt time.Time    `json:"created_at,omitempty"`
	//updatedAt time.Time    `json:"updated_at,omitempty"`
}

func (f *feedback) SetService(s ServiceValue) {
	f.Service = s
}

func (f *feedback) SetEdition(edition EditionValue) {
	f.Edition = edition
}

func (f *feedback) SetText(text TextValue) {
	f.Text = text
}

func (f *feedback) SetCustomer(customer Customer) {
	f.Customer = customer
}

func (f *feedback) SetContext(context Context) {
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
