package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request) error {
	_, _ = w.Write([]byte("Pong"))
	return nil
}

func store(w http.ResponseWriter, r *http.Request) error {
	// Reading the request body
	//body, err := ioutil.ReadAll(r.Body)
	body := make([]byte, env.MaxPostRequestBodyLength()+1)
	n, err := io.ReadFull(r.Body, body)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return errors.Wrap(err, "Error while reading request body")
		}
		body = body[:n]
	}
	if n < 1 {
		return makeBadRequestHttpError("Request body is empty", nil)
	}
	if uint(n) > env.MaxPostRequestBodyLength() {
		return makeRequestEntityTooLargeHttpError("Request body is too large", nil)
	}

	// Feedback
	feedback := &Feedback{}
	err = json.Unmarshal(body, feedback)
	if err != nil {
		return makeBadRequestHttpError("Request body is not a valid JSON string", err)
	}
	if feedback.Service == "" {
		return makeUnprocessableEntityHttpError("Property service is empty", nil)
	}
	if feedback.Text == "" {
		return makeUnprocessableEntityHttpError("Property text is empty", nil)
	}
	if feedback.Email != nil && *feedback.Email == "" {
		feedback.Email = nil
	}

	// Context
	feedback.Context = &Context{}
	err = json.Unmarshal(body, feedback.Context)
	if err != nil {
		return makeBadRequestHttpError("Error while unmarshalling feedback context from json",
			errors.Wrap(err, "Request body is not a valid JSON string"))
	}

	// Preparing context
	context, err := json.Marshal(feedback.Context)
	if err != nil {
		return errors.Wrap(err, "Error while marshalling context into json")
	}

	// Saving
	err = env.Db().QueryRow("INSERT INTO feedback (service, text, email, context) VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at",
		feedback.Id, feedback.Text, feedback.Email, context).Scan(&feedback.Id, &feedback.CreatedAt, &feedback.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "Error while saving feedback into database")
	}

	// Response
	json.NewEncoder(w).Encode(feedback)

	return nil
}
