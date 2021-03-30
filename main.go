package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mingalevme/feedbacker/infrastructure/config"
	"github.com/mingalevme/feedbacker/infrastructure/di"
	"github.com/mingalevme/feedbacker/infrastructure/http/router"
	"github.com/mingalevme/feedbacker/infrastructure/http/server"
	"github.com/pkg/errors"
	"os"
)

func main() {
	cfg := config.GetInstance()
	container := di.New(cfg)

	host := cfg.GetHTTPHost()
	port := cfg.GetHTTPPort()

	s := server.New(host, port, router.NewRouter(container))

	go func() {
		fmt.Printf("Listening on http://%s\n", s.Addr)
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
