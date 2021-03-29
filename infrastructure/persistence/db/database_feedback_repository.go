package db

import (
	"database/sql"
	"fmt"
	"github.com/mingalevme/feedbacker/domain/feedback"
	"github.com/pkg/errors"
	"time"
)

type FeedbackRepository struct {
	feedback.AbstractRepository
	db *sql.DB
}

func NewFeedbackRepository(connection *sql.DB) feedback.Repository {
	a := feedback.AbstractRepository{}
	return &FeedbackRepository{
		a,
		connection,
	}
}

type extraColumnHolder struct {
	Customer feedback.CustomerData `json:"customer"`
	Context feedback.ContextData   `json:"context"`
}

func (s *FeedbackRepository) CreateItem(service feedback.ServiceValue, edition feedback.EditionValue, text feedback.TextValue, customer feedback.Customer, context feedback.Context) (feedbackId feedback.FeedbackId, createdAt time.Time, err error) {
	extra := extraColumnHolder{
		Customer: customer.ToData(),
		Context: context.ToData(),
	}
	var id int
	query := "INSERT INTO feedback (service, edition, text, extra) VALUES ($1, $2, $3, $4) RETURNING id, created_at"
	err = s.db.QueryRow(query, service.Get(), edition.Get(), text.Get(), extra).Scan(&id, &createdAt)
	if err != nil {
		return nil, time.Now(), errors.Wrap(err, "Error while saving feedback into database")
	}
	feedbackId, err = feedback.NewFeedbackId(id)
	if err != nil {
		panic(fmt.Sprintf("Invalid id returning by database: %d", id))
	}
	return feedbackId, createdAt, nil
}
