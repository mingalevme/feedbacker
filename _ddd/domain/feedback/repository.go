// https://stackoverflow.com/questions/21250666/ddd-entity-identity-before-being-persisted

package feedback

import (
	"github.com/pkg/errors"
	"time"
)

var ErrNotFound = errors.New("model is not found")

// https://stackoverflow.com/a/39280987/1046909
type Repository interface {
	Add(service ServiceValue, edition EditionValue, text TextValue, customer Customer, context Context) (Feedback, error)
	//CreateItem(service ServiceValue, edition EditionValue, text TextValue, customer Customer, context Context) (FeedbackId, time.Time, error)
	GetById(id FeedbackID) (Feedback, error)
}

type AbstractRepository struct {
	RepositoryHelper
}

type RepositoryHelper interface {
	Repository
	//Add(service ServiceValue, edition EditionValue, text TextValue, customer Customer, context Context) (Feedback, error)
	CreateItem(service ServiceValue, edition EditionValue, text TextValue, customer Customer, context Context) (FeedbackID, time.Time, error)
	//GetById(id FeedbackId) (Feedback, error)
	MakeFeedback(id interface{}, service interface{}, edition interface{}, text interface{}, customer Customer, context Context, createdAt time.Time, updatedAt time.Time) Feedback
}

func (s *AbstractRepository) Add(service ServiceValue, edition EditionValue, text TextValue, customer Customer, context Context) (Feedback, error) {
	id, createdAt, err := s.RepositoryHelper.CreateItem(service, edition, text, customer, context)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating feedback item")
	}
	return s.MakeFeedback(id, service, edition, text, customer, context, createdAt, createdAt), nil
}

func (s *AbstractRepository) MakeFeedback(id interface{}, service interface{}, edition interface{}, text interface{}, customer Customer, context Context, createdAt time.Time, updatedAt time.Time) Feedback {
	fID, err := NewFeedbackID(id)
	if err != nil {
		panic("Error while creating FeedbackID from repository storage value")
	}
	serviceValue, err := NewServiceValue(service)
	if err != nil {
		panic("Error while creating ServiceValue from repository storage value")
	}
	editionValue, err := NewEditionValue(edition)
	if err != nil {
		panic("Error while creating EditionValue from repository storage value")
	}
	textValue, err := NewTextValue(text)
	if err != nil {
		panic("Error while creating TextValue from repository storage value")
	}
	data := FeedbackData{
		Id:        fID,
		Service:   serviceValue,
		Edition:   editionValue,
		Text:      textValue,
		Customer:  customer,
		Context:   context,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return &feedback{
		data,
	}
}

func (s *AbstractRepository) CreateItem(service ServiceValue, edition EditionValue, text TextValue, customer Customer, context Context) (FeedbackID, time.Time, error) {
	panic("not implemented")
}
