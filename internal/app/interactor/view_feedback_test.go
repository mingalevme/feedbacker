package interactor

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/internal/app/repository"
	"github.com/mingalevme/feedbacker/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestViewFeedbackSuccess(t *testing.T) {
	env := test.NewEnv(map[string]string{
		"PERSISTENCE_DRIVER": "array",
	})
	i := New(env)
	r, _ := env.FeedbackRepository().(*repository.ArrayFeedbackRepository)
	f1 := model.MakeFeedback()
	r.Storage = append(r.Storage, f1)
	f2, err := i.ViewFeedback(f1.ID)
	assert.NoError(t, err)
	assert.Equal(t, f1.ID, f2.ID)
}

func TestViewFeedbackNotFound(t *testing.T) {
	env := test.NewEnv(map[string]string{
		"PERSISTENCE_DRIVER": "array",
	})
	i := New(env)
	_, err := i.ViewFeedback(1)
	assert.ErrorIs(t, ErrNotFound, err)
}

