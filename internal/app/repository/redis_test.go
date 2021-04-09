package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/pkg/timeutils"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestRedisAddNoError(t *testing.T) {
	now := timeutils.Now()
	timeutils.SetTestNow(now)
	defer timeutils.ResetTestNow()

	client, mock := redismock.NewClientMock()

	r := NewRedisFeedbackRepository(client, context.Background())

	id := rand.Intn(1000) + 1
	mock.ExpectIncr(r.key + "_seq").SetVal(int64(id))

	f1 := model.MakeFeedback()
	data := AddFeedbackData{f1}

	f1.ID = id
	f1.CreatedAt = timeutils.Now()
	f1.UpdatedAt = timeutils.Now()

	j, err := json.Marshal(f1)
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectZAdd(r.key, &redis.Z{
		Score:  float64(id),
		Member: j,
	}).SetVal(1)

	f2, err := r.Add(data)
	assert.NoError(t, err)
	assert.Equal(t, id, f2.ID)
}

func TestRedisGetNoError(t *testing.T) {
	client, mock := redismock.NewClientMock()
	r := NewRedisFeedbackRepository(client, context.Background())
	f1 := model.MakeFeedback()
	j, err := json.Marshal(f1)
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectZRangeByScore(r.key, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%d", f1.ID),
		Max:    fmt.Sprintf("%d", f1.ID),
	}).SetVal([]string{string(j)})
	f2, err := r.GetById(f1.ID)
	assert.NoError(t, err)
	assert.Equal(t, f1.ID, f2.ID)
}

func TestRedisGetErrNotFound(t *testing.T) {
	client, mock := redismock.NewClientMock()
	r := NewRedisFeedbackRepository(client, context.Background())
	id := model.MakeFeedback().ID
	mock.ExpectZRangeByScore(r.key, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%d", id),
		Max:    fmt.Sprintf("%d", id),
	}).SetVal([]string{})
	_, err := r.GetById(id)
	assert.ErrorIs(t, err, ErrNotFound)
}
