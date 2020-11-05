package errors

import (
	jujuerr "github.com/juju/errors"
)

type databaseError struct {
	jujuerr.Err
}

// NewDatabaseErrorf returns an error which satisfies IsDatabaseError().
func NewDatabaseErrorf(format string, args ...interface{}) error {
	e := &databaseError{wrap(nil, format, " database error", args...)}
	e.SetLocation(1)
	return e
}

// NewDatabaseError returns an error which wraps err and satisfies IsDatabaseError().
func NewDatabaseError(err error, msg string) error {
	e := &databaseError{wrap(err, msg, " database error")}
	e.SetLocation(1)
	return e
}

// IsDatabaseError reports whether the error was created with
// NewDatabaseErrorf() or NewDatabaseError().
func IsDatabaseError(err error) bool {
	err = Cause(err)
	_, ok := err.(*databaseError)
	return ok
}
