package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var env = NewEnv()

func main() {

	host := env.GetEnv("FEEDBACKER_HOST", "")
	port := env.GetEnv("FEEDBACKER_PORT", "8080")

	r := mux.NewRouter()

	r.Handle("/ping", HttpRequestHandler{ping}).Methods("GET")
	r.Handle("/feedback", HttpRequestHandler{store}).Methods("POST")

	server := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", host, port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		env.Logger().Infof("Starting on http://%s:%s", host, port)
		if err := server.ListenAndServe(); err != nil {
			env.Logger().Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(server)

}

func waitForShutdown(server *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	server.Shutdown(ctx)

	env.Logger().Info("Shutting down")
	os.Exit(0)
}
