package main

import (
	_ "github.com/lib/pq"
	"github.com/mingalevme/feedbacker/infrastructure"
	"github.com/pkg/errors"
)

var env = NewEnv()

func main() {

	server := InitializeHttpServer()

	go func() {
		infrastructure.GetLogger().Infof("Starting on http://%s", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			env.Logger().Fatal(errors.Wrap(err, "Error while listening and serving"))
		}
	}()

	// Graceful Shutdown
	server.WaitForShutdown()

}
