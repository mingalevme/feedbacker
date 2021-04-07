package feedback

import (
	"github.com/mingalevme/feedbacker/_ddd/domain/feedback"
	"github.com/mingalevme/feedbacker/internal/app/service/log"
	"time"
)

type arrayFeedbackRepository struct {
	*feedback.AbstractRepository
	storage []feedback.Feedback
	logger  log.Logger
}

func NewArrayFeedbackRepository(logger log.Logger) feedback.Repository {
	parent := &feedback.AbstractRepository{}
	r := &arrayFeedbackRepository{
		parent,
		[]feedback.Feedback{},
		logger,
	}
	parent.RepositoryHelper = r
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

func (s *arrayFeedbackRepository) CreateItem(service feedback.ServiceValue, edition feedback.EditionValue, text feedback.TextValue, customer feedback.Customer, context feedback.Context) (feedback.FeedbackID, time.Time, error) {
	if id, err := feedback.NewFeedbackID(s.getNextID()); err != nil {
		return nil, time.Time{}, err
	} else {
		s.logger.Debugf("Created item: #%v, %s, %s, %s, %#v, %+v", id.Get(), service.Get(), edition.Get(), text.Get(), customer, context)
		return id, time.Now(), nil
	}
}

func (s *arrayFeedbackRepository) getNextID() int {
	return len(s.storage) + 1
}

func (s *arrayFeedbackRepository) GetById(id feedback.FeedbackID) (feedback.Feedback, error) {
	for _, f := range s.storage {
		if f.GetID().Equals(id) {
			//return f, nil
			return s.MakeFeedback(f.GetID(), f.GetService(), f.GetEdition(), f.GetText(), f.GetCustomer(), f.GetContext(), f.GetCreatedAt(), f.GetUpdatedAt()), nil
		}
	}
	return nil, feedback.ErrNotFound
}


