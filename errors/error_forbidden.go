package errors

import (
	jujuerr "github.com/juju/errors"
)

type forbiddenError struct {
	jujuerr.Err
}

// NewForbiddenErrorf returns an error which satisfies IsForbiddenError().
func NewForbiddenErrorf(format string, args ...interface{}) error {
	e := &forbiddenError{wrap(nil, format, " forbidden error", args...)}
	e.SetLocation(1)
	return e
}

// NewForbiddenError returns an error which wraps err and satisfies IsForbiddenError().
func NewForbiddenError(err error, msg string) error {
	e := &forbiddenError{wrap(err, msg, " forbidden error")}
	e.SetLocation(1)
	return e
}

// IsForbiddenError reports whether the error was created with
// NewForbiddenErrorf() or NewForbiddenError().
func IsForbiddenError(err error) bool {
	err = Cause(err)
	_, ok := err.(*forbiddenError)
	return ok
}
