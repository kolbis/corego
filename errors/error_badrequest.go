package errors

import (
	jujuerr "github.com/juju/errors"
)

type badreqError struct {
	jujuerr.Err
}

// NewBadRequestErrorf returns an error which satisfies IsBadRequestError().
func NewBadRequestErrorf(format string, args ...interface{}) error {
	e := &badreqError{wrap(nil, format, " bad request error", args...)}
	e.SetLocation(1)
	return e
}

// NewBadRequestError returns an error which wraps err and satisfies IsBadRequestError().
func NewBadRequestError(err error, msg string) error {
	e := &badreqError{wrap(err, msg, " bad request error")}
	e.SetLocation(1)
	return e
}

// IsBadRequestError reports whether the error was created with
// NewBadRequestErrorf() or NewBadRequestError().
func IsBadRequestError(err error) bool {
	err = Cause(err)
	_, ok := err.(*badreqError)
	return ok
}
