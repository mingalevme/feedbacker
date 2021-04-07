package repository

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/internal/app/service/log"
	"github.com/pkg/errors"
	"time"
)

type DatabaseFeedbackRepository struct {
	db     *sql.DB
	logger log.Logger
}

func NewDatabaseFeedbackRepository(conn *sql.DB, logger log.Logger) Feedback {
	return &DatabaseFeedbackRepository{
		conn,
		logger,
	}
}

type extraColumnHolder struct {
	Customer *model.Customer `json:"customer"`
	Context  *model.Context  `json:"context"`
}

func (c *extraColumnHolder) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (c *extraColumnHolder) Scan(src interface{}) error {
	var data []byte
	if b, ok := src.([]byte); ok {
		data = b
	} else if s, ok := src.(string); ok {
		data = []byte(s)
	}
	return json.Unmarshal(data, c)
}

func (s *DatabaseFeedbackRepository) Add(data AddFeedbackData) (model.Feedback, error) {
	if err := data.Validate(); err != nil {
		return model.Feedback{}, err
	}
	extra := &extraColumnHolder{
		Customer: data.Customer,
		Context:  data.Context,
	}
	if extra.Customer == nil && extra.Context == nil {
		extra = nil
	}
	var (
		id        int
		createdAt time.Time
	)
	s.logger.WithField("data", data).Debugf("Inserting feedback into database")
	query := "INSERT INTO feedback (service, edition, text, extra) VALUES ($1, $2, $3, $4) RETURNING id, created_at"
	err := s.db.QueryRow(query, data.Service, data.Edition, data.Text, extra).Scan(&id, &createdAt)
	if err != nil {
		return model.Feedback{}, errors.Wrap(err, "database: inserting feedback")
	}
	f := data.Feedback
	f.ID = id
	f.CreatedAt = createdAt
	f.UpdatedAt = createdAt
	return f, nil
}

func (s *DatabaseFeedbackRepository) GetById(id int) (model.Feedback, error) {
	var (
		service   string
		edition   string
		text      string
		createdAt time.Time
		updatedAt time.Time
		extra     extraColumnHolder
	)
	s.logger.WithField("id", id).Debugf("Selecting feedback from database")
	query := "SELECT service, edition, text, extra, created_at, updated_at FROM feedback WHERE id = $1"
	err := s.db.QueryRow(query, id).Scan(&service, &edition, &text, &extra, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return model.Feedback{}, ErrNotFound
	} else if err != nil {
		s.logger.WithError(err).WithField("id", id).Errorf("database: selecting feedback")
		return model.Feedback{}, ErrNotFound
	}
	f := model.Feedback{
		ID:        id,
		Service:   service,
		Edition:   &edition,
		Text:      text,
		Context:   nil,
		Customer:  nil,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	if extra.Customer != nil {
		f.Customer = extra.Customer
	}
	if extra.Context != nil {
		f.Context = extra.Context
	}
	return f, nil
}
