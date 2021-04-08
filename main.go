package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mingalevme/feedbacker/internal/app/di"
	"github.com/mingalevme/feedbacker/internal/config"
	"github.com/mingalevme/feedbacker/internal/http"
	"github.com/mingalevme/feedbacker/pkg/env"
	"github.com/pkg/errors"
)

func main() {
	e := env.New()
	c := config.New(e)
	container := di.New(c)
	address := e.Get("HTTP_LISTEN_ADDRESS", "0.0.0.0:8080")
	var s http.Server = http.NewEchoServer(address, container)
	go func() {
		fmt.Printf("Listening on http://%s\n", s.GetAddr())
		if err := s.ListenAndServe(); err != nil {
			panic(errors.Wrap(err, "http server: listening and serving"))
		}
	}()
	s.WaitForShutdown()
}
