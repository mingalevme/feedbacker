package repository

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/pkg/log"
	"time"
)

type ArrayFeedbackRepository struct {
	Storage []model.Feedback
	Logger  log.Logger
}

func NewArrayFeedbackRepository(logger log.Logger) *ArrayFeedbackRepository {
	return &ArrayFeedbackRepository{
		Storage: []model.Feedback{},
		Logger:  logger,
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
	s.Storage = append(s.Storage, f)
	s.Logger.Debugf("Created item: #%+v", f)
	return f, nil
}

func (s *ArrayFeedbackRepository) getNextID() int {
	return len(s.Storage) + 1
}

func (s *ArrayFeedbackRepository) GetById(id int) (model.Feedback, error) {
	for _, f := range s.Storage {
		if f.ID == id {
			return f, nil
		}
	}
	return model.Feedback{}, ErrNotFound
}