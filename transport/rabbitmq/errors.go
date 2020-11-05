package rabbitmq

import (
	"context"

	jujuerr "github.com/juju/errors"
	tlectx "github.com/kolbis/corego/context"
	"github.com/kolbis/corego/errors"
)

const (
	// NackErrorCode ...
	NackErrorCode int = 666
)

// Error ..
type Error struct {
	jujuerr.Err
	ErrorCode int
}

// NewRabbitErrorf returns an error which satisfies IsRabbitError().
func NewRabbitErrorf(code int, format string, args ...interface{}) error {
	newErr := jujuerr.NewErrWithCause(nil, format+" rabbitmq error", args...)
	e := &Error{
		newErr,
		code,
	}
	e.SetLocation(1)
	return e
}

// IsRabbitError reports whether the error was created with NewRabbitErrorf()
func IsRabbitError(err error) bool {
	err = errors.Cause(err)
	_, ok := err.(*Error)
	return ok
}

// SetNack will set in the context indication not to acknoledge the rabbit message
func SetNack(ctx context.Context) context.Context {
	err := NewRabbitErrorf(NackErrorCode, "Nack %d", NackErrorCode)
	return tlectx.SetError(ctx, err, NackErrorCode)
}

// ShouldNack will check if a proper Nack error was set into the context
// if it did, it will return true so that we will Nack the message
func ShouldNack(ctx context.Context) bool {
	hasError, code, err := tlectx.GetError(ctx)

	if hasError == true && IsRabbitError(err) && code == NackErrorCode {
		return true
	}

	return false
}
