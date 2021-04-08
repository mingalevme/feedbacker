package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/pkg/errors"
	"time"
)

type RedisFeedbackRepository struct {
	redis *redis.Client
	context context.Context
}

func NewRedisFeedbackRepository(redis *redis.Client, context context.Context) *RedisFeedbackRepository {
	return &RedisFeedbackRepository{
		redis: redis,
		context: context,
	}
}

func (s *RedisFeedbackRepository) Add(data AddFeedbackData) (model.Feedback, error)  {
	if err := data.Validate(); err != nil {
		return model.Feedback{}, err
	}
	f := data.Feedback
	f.ID = s.getNextID()
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	j, err := json.Marshal(f)
	if err != nil {
		return model.Feedback{}, errors.Wrap(err, "redis feedback repository: json-marshalling feedback")
	}
	_, err = s.redis.ZAdd(s.context, "feedbacks", &redis.Z{
		Score:  float64(f.ID),
		Member: j,
	}).Result()
	if err != nil {
		return model.Feedback{}, errors.Wrap(err, "redis feedback repository: add: z-adding feedback")
	}
	return f, nil
}

func (s *RedisFeedbackRepository) getNextID() int {
	id, err := s.redis.Incr(s.context, "feedback_seq").Result()
	if err != nil {
		panic(errors.Wrap(err, "redis feedback repository: get next id: incrementing feedback_seq"))
	}
	if id < 1 {
		id = 1
	}
	return int(id + 1)
}

func (s *RedisFeedbackRepository) GetById(id int) (model.Feedback, error) {
	z, err := s.redis.ZRangeByScore(s.context, "feedbacks", &redis.ZRangeBy{
		Min:    fmt.Sprintf("%d", id),
		Max:    fmt.Sprintf("%d", id),
	}).Result()
	if err != nil {
		return model.Feedback{}, errors.Wrap(err, "redis feedback repository: get: z-range-by-score")
	}
	if len(z) == 0 {
		return model.Feedback{}, ErrNotFound
	}
	f := &model.Feedback{}
	if err = json.Unmarshal([]byte(z[0]), f); err != nil {
		return model.Feedback{}, errors.Wrap(err, "redis feedback repository: json-unmarshalling feedback")
	}
	return *f, nil
}