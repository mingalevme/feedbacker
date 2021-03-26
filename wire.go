//+build wireinject

// https://dev.ms/2020/06/dependency-injection-google-wire/

package main

import (
	"github.com/google/wire"
	"github.com/mingalevme/feedbacker/app/services"
	"github.com/mingalevme/feedbacker/domain/feedback"
	"github.com/mingalevme/feedbacker/infrastructure"
	"github.com/mingalevme/feedbacker/infrastructure/http"
	"github.com/mingalevme/feedbacker/infrastructure/persistence"
	"github.com/mingalevme/feedbacker/infrastructure/persistence/db"
)

func InitializeFeedbackRepository() (feedback.Repository, error) {
	wire.Build(persistence.GetFeedbackRepository)
	return nil, nil
}

func InitializeDbFeedbackRepository() (feedback.Repository, error) {
	panic(wire.Build(db.NewFeedbackRepository, db.NewDatabaseConnection))
}

func InitializeHttpServer() *http.Server {
	panic(wire.Build(http.NewServer, infrastructure.GetLogger))
}

func InitializeLeaveFeedbackHandler() http.LeaveFeedbackHandler {
	panic(wire.Build(http.NewLeaveFeedbackHandler, services.NewLeaveFeedbackService, persistence.GetFeedbackRepository))
}
