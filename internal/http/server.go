package http

import (
	"context"
	"os"
	"os/signal"
	"time"
)

type Server interface {
	GetAddr() string
	ListenAndServe() error
	WaitForShutdown()
}

func waitForShutdown(shutdowner func(ctx context.Context) error) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := shutdowner(ctx); err != nil {
		panic(err)
	}

	os.Exit(0)
}