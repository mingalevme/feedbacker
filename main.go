package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mingalevme/feedbacker/internal/app/di"
	"github.com/mingalevme/feedbacker/internal/config"
	"github.com/mingalevme/feedbacker/internal/http"
	"github.com/mingalevme/feedbacker/pkg/env"
	"github.com/pkg/errors"
	"os"
)

func main() {
	envBag := env.New()
	cfg := config.New(envBag)
	container := di.New(cfg)

	address := envBag.Get("HTTP_LISTEN_ADDRESS", "0.0.0.0:8080")

	s := http.NewEchoServer(address, container)

	go func() {
		fmt.Printf("Listening on http://%s\n", s.GetAddr())
		if err1 := s.ListenAndServe(); err1 != nil {
			_, err2 := fmt.Fprintf(os.Stderr, "Error while listening and serving: %v\n", err1)
			if err2 != nil {
				panic(errors.Wrap(errors.Wrap(err1, "Error while listening and serving"), "Error stderr-ing error"))
			}
		}
	}()

	// Graceful Shutdown
	s.WaitForShutdown()
}
