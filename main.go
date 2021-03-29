package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mingalevme/feedbacker/infrastructure/cfg"
	"github.com/mingalevme/feedbacker/infrastructure/di"
	"github.com/mingalevme/feedbacker/infrastructure/env"
	"github.com/mingalevme/feedbacker/infrastructure/http/router"
	"github.com/mingalevme/feedbacker/infrastructure/http/server"
	"github.com/pkg/errors"
	"os"
)

func main() {
	config := cfg.New(env.New())
	container := di.New(config)

	host := config.GetHTTPHost()
	port := config.GetHTTPPort()

	s := server.New(host, port, router.NewRouter(container))

	go func() {
		fmt.Printf("Start listening http://%s\n", s.Addr)
		if err1 := s.ListenAndServe(); err1 != nil {
			if _, err2 := fmt.Fprintf(os.Stderr, "Error while listening and serving: %v\n", err1); err2 != nil {
				panic(errors.Wrap(errors.Wrap(err1, "Error while listening and serving"), "Error stderr-ing error"))
			}
		}
	}()

	// Graceful Shutdown
	s.WaitForShutdown()
}
