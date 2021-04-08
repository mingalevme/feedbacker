// +build testing

package test

import (
	"github.com/mingalevme/feedbacker/internal/app"
	"github.com/mingalevme/feedbacker/pkg/envvarbag"
)

func NewEnv(values map[string]string) *app.Container {
	return app.NewEnv(envvarbag.New().With(values))
}
