package main

import (
	"net/http"
)

type HttpError struct {
	statusCode int
	message    string
	cause      error
}

func (e HttpError) StatusCode() int {
	return e.statusCode
}

func (e HttpError) Cause() error {
	return e.cause
}

func (e HttpError) Message() string {
	return e.message
}

// Implements error interface
func (e HttpError) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

type HttpRequestHandler struct {
	handler func(w http.ResponseWriter, r *http.Request) error
}

func (h HttpRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.handler(w, r)
	if err != nil {
		switch e := err.(type) {
		case HttpError:
			http.Error(w, e.Error(), e.StatusCode())
		default:
			env.Logger().Errorf("%+v\n", e)
			if env.Debug() {
				http.Error(w, e.Error(), http.StatusInternalServerError)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
	}
}

func makeUnprocessableEntityHttpError(m string, err error) HttpError {
	return HttpError{
		statusCode: http.StatusUnprocessableEntity,
		message:    m,
		cause:      err,
	}
}

func makeBadRequestHttpError(m string, err error) HttpError {
	return HttpError{
		statusCode: http.StatusBadRequest,
		message:    m,
		cause:      err,
	}
}

func makeRequestEntityTooLargeHttpError(m string, err error) HttpError {
	return HttpError{
		statusCode: http.StatusRequestEntityTooLarge,
		message:    m,
		cause:      err,
	}
}
