package http

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
