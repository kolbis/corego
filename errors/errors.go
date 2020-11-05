package errors

import (
	jujuerr "github.com/juju/errors"
)

// New will create a new application error
func New(msg string) error {
	return newApplicationError(msg, 2)
}

// Annotate is used to add extra context to an existing error.
// The location of the Annotate call is recorded with the annotations.
// The file, line and function are also recorded.
func Annotate(err error, msg string) error {
	return jujuerr.Annotate(err, msg)
}

// Annotatef is used to add extra context to an existing error.
// The location of the Annotate call is recorded with the annotations.
// The file, line and function are also recorded.
func Annotatef(err error, format string, args ...interface{}) error {
	return jujuerr.Annotatef(err, format, args...)
}

// Cause returns the cause of the given error.
// This will be either the original error, or the result of a Wrap or Mask call.
// Cause is the usual way to diagnose errors that may have been wrapped by the other errors functions.
func Cause(err error) error {
	return jujuerr.Cause(err)
}

// Details returns information about the stack of errors wrapped by err, in the format:
// This is a terse alternative to ErrorStack as it returns a single line.
func Details(err error) string {
	return jujuerr.Details(err)
}

// ErrorStack returns a string representation of the annotated error.
// If the error passed as the parameter is not an annotated error, the result is simply the result of the Error() method on that error.
// If the error is an annotated error, a multi-line string is returned where each line represents one entry in the annotation stack.
// The full filename from the call stack is used in the output.
func ErrorStack(err error) string {
	return jujuerr.ErrorStack(err)
}

// Errorf creates a new annotated error and records the location that the error is created.
// This should be a drop in replacement for fmt.Errorf.
func Errorf(format string, args ...interface{}) error {
	return jujuerr.Errorf(format, args...)
}

// Mask hides the underlying error type, and records the location of the masking.
func Mask(err error) error {
	return jujuerr.Mask(err)
}

// Maskf masks the given error with the given format string and arguments (like fmt.Sprintf), returning a new error that maintains the error stack, but hides the underlying error type.
// The error string still contains the full annotations.
// If you want to hide the annotations, call Wrap.
func Maskf(err error, format string, args ...interface{}) error {
	return jujuerr.Maskf(err, format, args...)
}

// Wrap changes the Cause of the error.
// The location of the Wrap call is also stored in the error stack.
func Wrap(err error, newDescriptive error) error {
	return jujuerr.Wrap(err, newDescriptive)
}

// Wrapf changes the Cause of the error, and adds an annotation.
// The location of the Wrap call is also stored in the error stack
func Wrapf(other error, newDescriptive error, format string, args ...interface{}) error {
	return jujuerr.Wrapf(other, newDescriptive, format, args...)
}

// Trace adds the location of the Trace call to the stack.
// The Cause of the resulting error is the same as the error parameter.
// If the other error is nil, the result will be nil.
func Trace(err error) error {
	return jujuerr.Trace(err)
}

// wrap is a helper to construct an *wrapper.
func wrap(err error, format, suffix string, args ...interface{}) jujuerr.Err {
	newErr := jujuerr.NewErrWithCause(err, format+suffix, args...)
	newErr.SetLocation(2)
	return newErr
}
