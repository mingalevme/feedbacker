package service

import (
	"github.com/mingalevme/feedbacker/_ddd/domain/feedback"
	"github.com/mingalevme/feedbacker/internal/app/service/log"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.Errorf("%v", feedback.ErrNotFound)
//var ErrNotFound = fmt.Errorf("%v", feedback.ErrNotFound)

var ErrUnprocessableEntity = errors.New("unprocessable entity")

type ViewFeedbackData struct {
	Id interface{} `json:"id,omitempty"`
}

type ViewFeedbackService interface {
	Handle(id ViewFeedbackData) (feedback.FeedbackConstData, error)
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

func (s *viewFeedbackService) Handle(data ViewFeedbackData) (feedback.FeedbackConstData, error) {
	id, err := feedback.NewFeedbackID(data.Id)
	if err != nil {
		return nil, errors.Wrap(ErrUnprocessableEntity, err.Error())
	}
	f, err := s.repository.GetById(id)
	if errors.Is(err, feedback.ErrNotFound) {
		return nil, ErrNotFound
	}
	return feedback.NewFeedbackData(f), nil
}
