package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mingalevme/feedbacker/internal/app"
	"github.com/mingalevme/feedbacker/internal/http"
	"github.com/mingalevme/feedbacker/pkg/envvarbag"
	"github.com/pkg/errors"
)

func main() {
	envVarBag := envvarbag.New()
	var env app.Env = app.NewEnv(envVarBag)
	address := envVarBag.Get("HTTP_LISTEN_ADDRESS", "0.0.0.0:8080")
	var s http.Server = http.NewEchoServer(address, env)
	go func() {
		fmt.Printf("Listening on http://%s\n", s.GetAddr())
		if err := s.ListenAndServe(); err != nil {
			panic(errors.Wrap(err, "http server: listening and serving"))
		}
	}()
	s.WaitForShutdown()
}
