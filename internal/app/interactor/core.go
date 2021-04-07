package interactor

import (
	"github.com/mingalevme/feedbacker/internal/app/di"
	"github.com/mingalevme/feedbacker/internal/app/repository"
	"github.com/pkg/errors"
)

var ErrUnprocessableEntity = errors.New(repository.ErrUnprocessableEntity.Error())
var ErrNotFound = errors.New(repository.ErrNotFound.Error())

type Interactor struct {
	container di.Container
}

func New(container di.Container) *Interactor {
	return &Interactor{
		container: container,
	}
}