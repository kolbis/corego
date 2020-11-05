package errors

import (
	jujuerr "github.com/juju/errors"
)

type notfoundError struct {
	jujuerr.Err
}

// NewNotFoundErrorf returns an error which satisfies IsNotFoundError().
func NewNotFoundErrorf(format string, args ...interface{}) error {
	e := &notfoundError{wrap(nil, format, " notfound error", args...)}
	e.SetLocation(1)
	return e
}

// NewNotFoundError returns an error which wraps err and satisfies IsNotFoundError().
func NewNotFoundError(err error, msg string) error {
	e := &notfoundError{wrap(err, msg, " notfound error")}
	e.SetLocation(1)
	return e
}

// IsNotFoundError reports whether the error was created with
// NewNotFoundError() or NewApplicationErrorf().
func IsNotFoundError(err error) bool {
	err = Cause(err)
	_, ok := err.(*notfoundError)
	return ok
}
