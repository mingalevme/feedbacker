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
	CreateItem(service ServiceValue, edition EditionValue, text TextValue, customer Customer, context Context) (FeedbackId, time.Time, error)
	GetById(id FeedbackId) (Feedback, error)
}

type AbstractRepository struct {
	Repository
}

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
