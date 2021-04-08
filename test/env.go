// +build testing

package test

import (
	"github.com/mingalevme/feedbacker/internal/app"
	"github.com/mingalevme/feedbacker/pkg/envvarbag"
)

func NewEnv(values map[string]string) *app.Container {
	if values == nil {
		values = map[string]string{}
	}
	return app.NewEnv(envvarbag.NewWithValues(values))
}
