package errors

import (
	jujuerr "github.com/juju/errors"
)

type timeoutError struct {
	jujuerr.Err
}

// NewTimeoutErrorf returns an error which satisfies IsTimeoutError().
func NewTimeoutErrorf(format string, args ...interface{}) error {
	e := &timeoutError{wrap(nil, format, " application error", args...)}
	e.SetLocation(1)
	return e
}

// NewTimeoutError returns an error which wraps err and satisfies IsTimeoutError().
func NewTimeoutError(err error, msg string) error {
	e := &timeoutError{wrap(err, msg, "")}
	e.SetLocation(1)
	return e
}

// IsTimeoutError reports whether the error was created with
// NewTimeoutError() or NewTimeoutErrorf().
func IsTimeoutError(err error) bool {
	err = Cause(err)
	_, ok := err.(*timeoutError)
	return ok
}
