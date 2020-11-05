package errors

import (
	jujuerr "github.com/juju/errors"
)

type methodNotAllowedError struct {
	jujuerr.Err
}

// NewMethodNotAllowedf returns an error which satisfies IsMethodNotAllowedError().
func NewMethodNotAllowedf(format string, args ...interface{}) error {
	e := &methodNotAllowedError{wrap(nil, format, " method not allowed error", args...)}
	e.SetLocation(1)
	return e
}

// NewMethodNotAllowedError returns an error which wraps err and satisfies IsMethodNotAllowedError().
func NewMethodNotAllowedError(err error, msg string) error {
	e := &methodNotAllowedError{wrap(err, msg, " method not allowed error")}
	e.SetLocation(1)
	return e
}

// IsMethodNotAllowedError reports whether the error was created with
// NewMethodNotAllowedError() or NewMethodNotAllowedf().
func IsMethodNotAllowedError(err error) bool {
	err = Cause(err)
	_, ok := err.(*methodNotAllowedError)
	return ok
}
