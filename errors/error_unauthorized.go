package errors

import (
	jujuerr "github.com/juju/errors"
)

type unauthorizedError struct {
	jujuerr.Err
}

// NewUnauthorizedErrorf returns an error which satisfies IsUnauthorizedError().
func NewUnauthorizedErrorf(format string, args ...interface{}) error {
	e := &unauthorizedError{wrap(nil, format, " unauthorized error", args...)}
	e.SetLocation(1)
	return e
}

// NewUnauthorizedError returns an error which wraps err and satisfies IsUnauthorizedError().
func NewUnauthorizedError(err error, msg string) error {
	e := &unauthorizedError{wrap(err, msg, " unauthorized error")}
	e.SetLocation(1)
	return e
}

// IsUnauthorizedError reports whether the error was created with
// NewUnauthorizedErrorf() or NewUnauthorizedError().
func IsUnauthorizedError(err error) bool {
	err = Cause(err)
	_, ok := err.(*unauthorizedError)
	return ok
}
