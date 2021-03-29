package persistence

import (
	"github.com/mingalevme/feedbacker/domain/feedback"
	"github.com/mingalevme/feedbacker/infrastructure/log"
	"time"
)

type loggerFeedbackRepository struct {
	feedback.AbstractRepository
	logger log.Logger
}

func NewLoggerFeedbackRepository(logger log.Logger) feedback.Repository {
	a := feedback.AbstractRepository{}
	return &loggerFeedbackRepository{
		a,
		logger,
	}
}

func (s *loggerFeedbackRepository) CreateItem(service feedback.ServiceValue, edition feedback.EditionValue, text feedback.TextValue, customer feedback.Customer, context feedback.Context) (feedbackId feedback.FeedbackId, createdAt time.Time, err error) {
	s.logger.Debugf("Creating item: %s, %s, %s, %s, %s", service.Get(), edition.Get(), text.Get(), customer, context)
	id, _ := feedback.NewFeedbackId(1)
	return id, time.Now(), nil
}
