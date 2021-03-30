package persistence

import (
	"github.com/mingalevme/feedbacker/domain/feedback"
	"github.com/mingalevme/feedbacker/infrastructure/log"
	"time"
)

type arrayFeedbackRepository struct {
	*feedback.AbstractRepository
	storage []feedback.Feedback
	logger log.Logger
}

func NewArrayFeedbackRepository(logger log.Logger) feedback.Repository {
	parent := &feedback.AbstractRepository{}
	r := &arrayFeedbackRepository{
		parent,
		[]feedback.Feedback{},
		logger,
	}
	parent.Repository = r
	return r
}

func (s *arrayFeedbackRepository) Add(service feedback.ServiceValue, edition feedback.EditionValue, text feedback.TextValue, customer feedback.Customer, context feedback.Context) (feedback.Feedback, error)  {
	f, err := s.AbstractRepository.Add(service, edition, text, customer, context)
	if err != nil {
		return nil, err
	}
	s.storage = append(s.storage, f)
	return f, nil
}

var i int

func (s *arrayFeedbackRepository) CreateItem(service feedback.ServiceValue, edition feedback.EditionValue, text feedback.TextValue, customer feedback.Customer, context feedback.Context) (feedback.FeedbackId, time.Time, error) {
	i++
	if id, err := feedback.NewFeedbackId(i); err != nil {
		return nil, time.Time{}, err
	} else {
		s.logger.Debugf("Created item: #%v, %s, %s, %s, %#v, %+v", id.Get(), service.Get(), edition.Get(), text.Get(), customer, context)
		return id, time.Now(), nil
	}
}

func (s *arrayFeedbackRepository) GetById(id feedback.FeedbackId) (feedback.Feedback, error) {
	for _, f := range s.storage {
		if f.GetID().Equals(id) {
			return f, nil
		}
	}
	return nil, feedback.ErrNotFound
}


