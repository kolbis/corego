package errors

import (
	jujuerr "github.com/juju/errors"
)

type internalServerError struct {
	jujuerr.Err
}

// NewInternalServerErrorf returns an error which satisfies IsInternalServerError().
func NewInternalServerErrorf(format string, args ...interface{}) error {
	e := &internalServerError{wrap(nil, format, " internal server error", args...)}
	e.SetLocation(1)
	return e
}

// NewInternalServerError returns an error which wraps err and satisfies IsInternalServerError().
func NewInternalServerError(err error, msg string) error {
	e := &internalServerError{wrap(err, msg, " internal server error")}
	e.SetLocation(1)
	return e
}

// IsInternalServerError reports whether the error was created with
// NewInternalServerErrorf() or NewInternalServerError().
func IsInternalServerError(err error) bool {
	err = Cause(err)
	_, ok := err.(*internalServerError)
	return ok
}
