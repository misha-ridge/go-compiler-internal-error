package apierror

import (
	"errors"
	"fmt"
	"net/http"
)

// Error is an extended error interface that can tell with which HTTP status to
// report it.
//
// An error impementing this interface has a curated message that can be
// returned to the API client. Errors that don't implement this interface might
// leak sensitive information, so their messages must not be returned to the API
// client.
type Error interface {
	error
	HTTPStatus() int
}

type wrapped struct {
	error
	status int
}

func (w wrapped) HTTPStatus() int {
	return w.status
}

func (w wrapped) Unwrap() error {
	return errors.Unwrap(w.error)
}

// Wrap adds an HTTP status to an ordinary error;
// returns nil if err is nil.
//
// Note that wrapping an error blesses its message as safe to return to the API
// client; be careful about what you wrap.
func Wrap(status int, err error) Error {
	if err == nil {
		return nil
	}
	return wrapped{error: err, status: status}
}

// Errorf adds an HTTP status to an error returned by fmt.Errorf
func Errorf(status int, format string, a ...any) Error {
	return Wrap(status, fmt.Errorf(format, a...))
}

// BadRequest returns a 400 Bad Request error with the given message
func BadRequest(msg string, a ...any) Error {
	return Errorf(http.StatusBadRequest, msg, a...)
}

// NotFound returns a 404 Not Found error with the given message
func NotFound(msg string, a ...any) Error {
	return Errorf(http.StatusNotFound, msg, a...)
}

// AccessDenied is the most commonly used 403 Forbidden error
var AccessDenied = Errorf(http.StatusForbidden, "access denied")
