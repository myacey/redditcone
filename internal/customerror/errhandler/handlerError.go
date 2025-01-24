package errhandler

import (
	"errors"
	"net/http"
)

// StatusCodedError represents an error with an HTTP status code,
// a user message, and additional information for logging.
type StatusCodedError struct {
	StatusCode int    `json:"-"`
	UserAnswer string `json:"message"`
	DebugLog   string `json:"-"`
	Err        error  `json:"-"`
}

// Error implements the error interface.
// Returns UserAnswer
func (e *StatusCodedError) Error() string {
	if e.UserAnswer != "" {
		return e.UserAnswer
	}
	if e.DebugLog != "" {
		return e.DebugLog
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "unknown error"
}

// NewStatusCodedError creates a new instance of StatusCodedError
func New(statusCode int, userAnswer, debugLog string, err error) error {
	return &StatusCodedError{
		StatusCode: statusCode,
		UserAnswer: userAnswer,
		DebugLog:   debugLog,
		Err:        err,
	}
}

// GetStatusCode returns the HTTP status code from the error.
// Default is BadRequest
func GetStatusCode(err error) int {
	var scErr *StatusCodedError
	if errors.As(err, &scErr) {
		return scErr.StatusCode
	}
	return http.StatusBadRequest // Default status code
}
