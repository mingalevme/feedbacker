package service

import (
	"github.com/mingalevme/feedbacker/domain/feedback"
	"github.com/mingalevme/feedbacker/infrastructure/log"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.Errorf("%v", feedback.ErrNotFound)
//var ErrNotFound = fmt.Errorf("%v", feedback.ErrNotFound)

var ErrUnprocessableEntity = errors.New("unprocessable entity")

type ViewFeedbackData struct {
	Id interface{} `json:"id,omitempty"`
}

type ViewFeedbackService interface {
	Handle(id ViewFeedbackData) (feedback.FeedbackData, error)
}

type viewFeedbackService struct {
	repository feedback.Repository
	logger     log.Logger
}

func NewViewFeedbackService(repository feedback.Repository, logger log.Logger) ViewFeedbackService {
	if repository == nil {
		panic(errors.New("repository is nil"))
	}
	if logger == nil {
		panic(errors.New("logger is nil"))
	}
	service := viewFeedbackService{
		repository: repository,
		logger:     logger,
	}
	return &service
}

func (s *viewFeedbackService) Handle(data ViewFeedbackData) (feedback.FeedbackData, error) {
	id, err := feedback.NewFeedbackId(data.Id)
	if err != nil {
		return nil, errors.Wrap(ErrUnprocessableEntity, "id is invalid")
	}
	f, err := s.repository.GetById(id)
	if errors.Is(err, feedback.ErrNotFound) {
		return nil, ErrNotFound
	}
	return feedback.NewFeedbackData(f), nil
}
