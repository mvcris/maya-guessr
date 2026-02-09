package coreerrors

import (
	"errors"
	"fmt"
)

// HTTPStatusCoder is implemented by domain errors that carry an HTTP status code.
type HTTPStatusCoder interface {
	error
	HTTPStatus() int
}

// domainError holds a message and HTTP status code.
type domainError struct {
	message string
	status  int
}

func (e *domainError) Error() string {
	return e.message
}

func (e *domainError) HTTPStatus() int {
	return e.status
}

// Conflict returns an error with HTTP status 409.
func Conflict(message string) HTTPStatusCoder {
	return &domainError{message: message, status: 409}
}

// Forbidden returns an error with HTTP status 403.
func Forbidden(message string) HTTPStatusCoder {
	return &domainError{message: message, status: 403}
}

// NotFound returns an error with HTTP status 404.
func NotFound(message string) HTTPStatusCoder {
	return &domainError{message: message, status: 404}
}

// Unauthorized returns an error with HTTP status 401.
func Unauthorized(message string) HTTPStatusCoder {
	return &domainError{message: message, status: 401}
}

// BadRequest returns an error with HTTP status 400.
func BadRequest(message string) HTTPStatusCoder {
	return &domainError{message: message, status: 400}
}

// Validation returns an error with HTTP status 400 for validation failures.
// If details is nil or empty, message is used as the single error message.
func Validation(message string, details map[string]string) HTTPStatusCoder {
	if len(details) == 0 {
		return &domainError{message: message, status: 400}
	}
	msg := message
	for k, v := range details {
		msg = fmt.Sprintf("%s; %s: %s", msg, k, v)
	}
	return &domainError{message: msg, status: 400}
}

// Status returns the HTTP status code for err if it is a domain error, and true.
// Otherwise it returns (0, false).
func Status(err error) (int, bool) {
	var e *domainError
	if errors.As(err, &e) {
		return e.status, true
	}
	return 0, false
}
