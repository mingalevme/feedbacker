package persistence

import (
	"github.com/mingalevme/feedbacker"
	"github.com/mingalevme/feedbacker/domain/feedback"
	"github.com/mingalevme/feedbacker/infrastructure/env"
	"github.com/mingalevme/feedbacker/infrastructure/persistence/db"
	"github.com/pkg/errors"
)

func GetFeedbackRepository() (feedback.Repository, error) {
	driver := env.GetEnvValue("FEEDBACKER_PERSISTENCE_DRIVER", "database")
	if driver == "database" {
		return db.NewFeedbackRepository(db.NewDatabaseConnection())
	} else {
		return nil, errors.Errorf("Unsupported persistence driver: %s", driver)
	}
}
