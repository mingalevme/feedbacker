// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/mingalevme/feedbacker/domain/feedback"
	"github.com/mingalevme/feedbacker/infrastructure"
	"github.com/mingalevme/feedbacker/infrastructure/http"
	"github.com/mingalevme/feedbacker/infrastructure/persistence/db"
)

// Injectors from wire.go:

func InitializeFeedbackRepository() (feedback.Repository, error) {
	repository, err := NewFeedbackRepository()
	if err != nil {
		return nil, err
	}
	return repository, nil
}

func InitializeDbFeedbackRepository() (feedback.Repository, error) {
	sqlDB, err := NewDatabaseConnection()
	if err != nil {
		return nil, err
	}
	repository := db.NewFeedbackRepository(sqlDB)
	return repository, nil
}

func InitializeHttpServer() *http.Server {
	logrusLogger := infrastructure.GetLogger()
	server := http.NewServer(logrusLogger)
	return server
}