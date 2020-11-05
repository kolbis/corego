package errors

import (
	jujuerr "github.com/juju/errors"
)

type applicationError struct {
	jujuerr.Err
}

// NewApplicationErrorf returns an error which satisfies IsNotSupported().
func NewApplicationErrorf(format string, args ...interface{}) error {
	e := &applicationError{wrap(nil, format, " application error", args...)}
	e.SetLocation(1)
	return e
}

func newApplicationError(msg string, location int) error {
	e := &applicationError{wrap(nil, msg, " application error")}
	e.SetLocation(location)
	return e
}

// NewApplicationError returns an error which wraps err and satisfies IsApplicationError().
func NewApplicationError(err error, msg string) error {
	e := &applicationError{wrap(err, msg, " application error")}
	e.SetLocation(1)
	return e
}

// IsApplicationError reports whether the error was created with
// NotSupportedf() or NewNotSupported().
func IsApplicationError(err error) bool {
	err = Cause(err)
	_, ok := err.(*applicationError)
	return ok
}
