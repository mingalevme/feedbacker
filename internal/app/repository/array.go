package repository

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/internal/app/service/log"
	"time"
)

type ArrayFeedbackRepository struct {
	storage []model.Feedback
	logger  log.Logger
}

func NewArrayFeedbackRepository(logger log.Logger) *ArrayFeedbackRepository {
	return &ArrayFeedbackRepository{
		storage: []model.Feedback{},
		logger: logger,
	}
}

func (s *ArrayFeedbackRepository) Add(data AddFeedbackData) (model.Feedback, error)  {
	if err := data.Validate(); err != nil {
		return model.Feedback{}, err
	}
	f := data.Feedback
	f.ID = s.getNextID()
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	s.storage = append(s.storage, f)
	s.logger.Debugf("Created item: #%+v", f)
	return f, nil
}

func (s *ArrayFeedbackRepository) getNextID() int {
	return len(s.storage) + 1
}

func (s *ArrayFeedbackRepository) GetById(id int) (model.Feedback, error) {
	for _, f := range s.storage {
		if f.ID == id {
			return f, nil
		}
	}
	return model.Feedback{}, ErrNotFound
}