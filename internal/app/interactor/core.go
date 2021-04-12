package interactor

import (
	"github.com/mingalevme/feedbacker/internal/app"
	"github.com/mingalevme/feedbacker/internal/app/repository"
	"github.com/pkg/errors"
	"sync"
)

var ErrUnprocessableEntity = errors.New(repository.ErrUnprocessableEntity.Error())
var ErrNotFound = errors.New(repository.ErrNotFound.Error())

type Interactor struct {
	env app.Env
	wg  sync.WaitGroup
}

func New(env app.Env) *Interactor {
	return &Interactor{
		env: env,
	}
}
