package http

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mingalevme/feedbacker/infrastructure/env"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	http.Server
	logger *log.Logger
}

func NewServer(logger *log.Logger) *Server {
	host := env.GetEnvValue("FEEDBACKER_HOST", "")
	port := env.GetEnvValue("FEEDBACKER_PORT", "8080")

	r := mux.NewRouter()

	r.Handle("/ping", NewRequestHandler(ping)).Methods("GET")
	//r.Handle("/feedback", NewLeaveFeedbackHandler(services.NewLeaveFeedbackService())).Methods("POST")

	server := &Server{
		http.Server{
			Handler:      r,
			Addr:         fmt.Sprintf("%s:%s", host, port),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		logger,
	}

	return server
}

func (s *Server) WaitForShutdown() {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_ = s.Shutdown(ctx)

	//env.Logger().Info("Shutting down")
	os.Exit(0)
}
