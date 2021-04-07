package feedback

import (
	"database/sql"
	"fmt"
	"github.com/mingalevme/feedbacker/_ddd/domain/feedback"
	"github.com/mingalevme/feedbacker/internal/app/service/log"
	"github.com/pkg/errors"
	"time"
)

type databaseFeedbackRepository struct {
	*feedback.AbstractRepository
	db     *sql.DB
	logger log.Logger
}

func NewFeedbackRepository(conn *sql.DB, logger log.Logger) feedback.Repository {
	return &databaseFeedbackRepository{
		&feedback.AbstractRepository{},
		conn,
		logger,
	}
}

type extraColumnHolder struct {
	Customer *feedback.CustomerData `json:"customer"`
	Context *feedback.ContextData   `json:"context"`
}

func (s *databaseFeedbackRepository) CreateItem(service feedback.ServiceValue, edition feedback.EditionValue, text feedback.TextValue, customer feedback.Customer, context feedback.Context) (feedback.FeedbackID, time.Time, error) {
	extra := &extraColumnHolder{
		Customer: feedback.GetCustomerMapper().CustomerToData(customer),
		Context: feedback.GetContextMapper().ContextToData(context),
	}
	if extra.Customer == nil && extra.Context == nil {
		extra = nil
	}
	var (
		id int
		createdAt time.Time
	)
	query := "INSERT INTO feedback (service, edition, text, extra) VALUES ($1, $2, $3, $4) RETURNING id, created_at"
	err := s.db.QueryRow(query, service.Get(), edition.Get(), text.Get(), extra).Scan(&id, &createdAt)
	if err != nil {
		return nil, time.Now(), errors.Wrap(err, "Error while saving feedback into database")
	}
	feedbackId, err := feedback.NewFeedbackID(id)
	if err != nil {
		panic(fmt.Sprintf("Invalid id returning by database: %d", id))
	}
	return feedbackId, createdAt, nil
}

func (s *databaseFeedbackRepository) GetById(id feedback.FeedbackID) (feedback.Feedback, error) {
	var (
		service   string
		edition   string
		text      string
		createdAt time.Time
		updatedAt time.Time
		extra     extraColumnHolder
	)
	query := "SELECT service, edition, text, extra, created_at, updated_at FROM feedback WHERE id = $1"
	err := s.db.QueryRow(query, id.Get()).Scan(&service, &edition, &text, &extra, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, feedback.ErrNotFound
	} else if err != nil {
		s.logger.WithError(err).WithField("feedback-id", id.Get()).Errorf("Error while selecting feedback from database")
		return nil, feedback.ErrNotFound
	}
	customer, err := feedback.GetCustomerMapper().DataToCustomer(extra.Customer)
	if err != nil {
		panic("Error while creating Customer from database row value")
	}
	context, err := feedback.GetContextMapper().DataToContext(extra.Context)
	if err != nil {
		panic("Error while creating Context from database row value")
	}
	return s.MakeFeedback(id, service, edition, text, customer, context, createdAt, updatedAt), nil
}
