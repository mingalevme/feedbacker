package responder

import (
	"encoding/json"
	"github.com/mingalevme/feedbacker/internal/app/service/log"
	"github.com/pkg/errors"
	"net/http"
)

type Responder interface {
	Respond(w http.ResponseWriter, v interface{}) error
}

type jsonResponder struct {
	logger log.Logger
}

func NewJSONResponder(logger log.Logger) Responder {
	return &jsonResponder{
		logger: logger,
	}
}

func (s *jsonResponder) Respond(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		err = errors.Wrap(err, "Error while writing json-response")
		panic(err)
		//s.logger.Error(err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//return err
	}
	return nil
}