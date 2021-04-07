package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	netHTTP "net/http"
)

type Handler netHTTP.Handler

func requestToData(r *netHTTP.Request, data interface{}) error {
	params := mux.Vars(r)
	j, err := json.Marshal(params)
	if err != nil {
		panic(errors.Wrap(err, "Error while marshaling to json"))
	}
	err = json.Unmarshal(j, data)
	if err != nil {
		return err
	}
	return nil
}
