package http

import (
	"context"
)

type Server interface {
	GetAddr() string
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}
