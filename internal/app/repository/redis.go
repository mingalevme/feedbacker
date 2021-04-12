package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/pkg/timeutils"
	"github.com/pkg/errors"
)

type RedisFeedbackRepository struct {
	redis   *redis.Client
	key     string
	context context.Context
}

func NewRedisFeedbackRepository(redis *redis.Client, context context.Context) *RedisFeedbackRepository {
	return &RedisFeedbackRepository{
		redis:   redis,
		key:     "feedback",
		context: context,
	}
}

func (s *RedisFeedbackRepository) Name() string {
	return "redis"
}

func (s *RedisFeedbackRepository) Add(data AddFeedbackData) (model.Feedback, error) {
	if err := data.Validate(); err != nil {
		return model.Feedback{}, err
	}
	f := data.Feedback
	f.ID = s.getNextID()
	f.CreatedAt = timeutils.Now()
	f.UpdatedAt = timeutils.Now()
	j, err := json.Marshal(f)
	if err != nil {
		return model.Feedback{}, errors.Wrap(err, "redis feedback repository: json-marshalling feedback")
	}
	_, err = s.redis.ZAdd(s.context, s.key, &redis.Z{
		Score:  float64(f.ID),
		Member: j,
	}).Result()
	if err != nil {
		return model.Feedback{}, errors.Wrap(err, "redis feedback repository: add: z-adding feedback")
	}
	return f, nil
}

func (s *RedisFeedbackRepository) getNextID() int {
	id, err := s.redis.Incr(s.context, fmt.Sprintf("%s_seq", s.key)).Result()
	if err != nil {
		panic(errors.Wrap(err, "redis feedback repository: get next id: incrementing feedback_seq"))
	}
	if id < 1 {
		id = 1
	}
	return int(id)
}

func (s *RedisFeedbackRepository) GetById(id int) (model.Feedback, error) {
	z, err := s.redis.ZRangeByScore(s.context, s.key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", id),
		Max: fmt.Sprintf("%d", id),
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

func (s *RedisFeedbackRepository) Health() error {
	return s.redis.Ping(s.context).Err()
}
