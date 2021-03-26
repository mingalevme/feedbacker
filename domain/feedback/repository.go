// https://stackoverflow.com/questions/21250666/ddd-entity-identity-before-being-persisted

package feedback

import (
	"github.com/pkg/errors"
	"time"
)

type Repository interface {
	Add(service ServiceValue, Edition EditionValue, text TextValue, customer Customer, context Context) (Feedback, error)
}

type AbstractRepository struct {}

func (s *AbstractRepository) Add(service ServiceValue, edition EditionValue, text TextValue, customer Customer, context Context) (Feedback, error) {
	id, createdAt, err := s.CreateItem(service, edition, text, customer, context)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating feedback item")
	}
	data := feedbackData{
		Id:        id,
		Service:   service,
		Edition:   edition,
		Text:      text,
		Context:   context,
		Customer:  customer,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
	return &feedback{
		data,
	}, nil
}

func (s *AbstractRepository) CreateItem(service ServiceValue, edition EditionValue, text TextValue, customer Customer, context Context) (feedbackId FeedbackId, createdAt time.Time, err error) {
	panic("Method is not implemented")
}

