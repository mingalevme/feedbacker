package feedback

import (
	"time"
)

type FeedbackData struct {
	Id        int           `json:"id,omitempty"`
	Service   string        `json:"service,omitempty"`
	Edition   string        `json:"edition,omitempty"`
	Text      string        `json:"text,omitempty"`
	Customer  *CustomerData `json:"customer,omitempty"`
	Context   *ContextData  `json:"context,omitempty"`
	CreatedAt time.Time     `json:"created_at,omitempty"`
	UpdatedAt time.Time     `json:"updated_at,omitempty"`
}

//type FeedbackConstData interface {
//	GetID() FeedbackID
//	GetService() ServiceValue
//	GetEdition() EditionValue
//	GetText() TextValue
//	GetCustomer() Customer
//	GetContext() Context
//	GetCreatedAt() time.Time
//	GetUpdatedAt() time.Time
//	Equals(value interface{}) bool
//	MarshalJSON() ([]byte, error)
//}

//func (f FeedbackData) GetID() FeedbackID {
//	return f.Id
//}
//
//func (f FeedbackData) GetService() ServiceValue {
//	return f.Service
//}
//
//func (f FeedbackData) GetEdition() EditionValue {
//	return f.Edition
//}
//
//func (f FeedbackData) GetText() TextValue {
//	return f.Text
//}
//
//func (f FeedbackData) GetCustomer() Customer {
//	return GetCustomerMapper().DataToCustomer(f.Customer)
//}
//
//func (f FeedbackData) GetContext() Context {
//	return f.Context
//}
//
//func (f FeedbackData) GetCreatedAt() time.Time {
//	return f.CreatedAt
//}
//
//func (f FeedbackData) GetUpdatedAt() time.Time {
//	return f.UpdatedAt
//}
//
//func (f FeedbackData) Equals(value interface{}) bool {
//	target, ok := value.(FeedbackConstData)
//	return ok && f.Id.Equals(target.GetID())
//}

//func (f FeedbackData) MarshalJSON() ([]byte, error) {
//	panic("not implemented")
//	//return json.Marshal(&struct {
//	//	ID       int64  `json:"id"`
//	//	Name     string `json:"name"`
//	//	LastSeen int64  `json:"lastSeen"`
//	//}{
//	//	ID:       u.ID,
//	//	Name:     u.Name,
//	//	LastSeen: u.LastSeen.Unix(),
//	//})
//}

//func (f FeedbackData) UnmarshalJSON(data []byte) error {
//	value, err := json.Unmarshal(...)
//	return nil
//}

//func NewFeedbackData(f Feedback) FeedbackConstData {
//	return FeedbackData{
//		Id:        f.GetID(),
//		Service:   f.GetService(),
//		Edition:   f.GetEdition(),
//		Text:      f.GetText(),
//		Customer:  f.GetCustomer(),
//		Context:   f.GetContext(),
//		CreatedAt: f.GetCreatedAt(),
//		UpdatedAt: f.GetUpdatedAt(),
//	}
//}
