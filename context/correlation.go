package context

import (
	"github.com/kolbis/corego/utils"
)

// Correlation used to create correlation ID
type Correlation interface {
	New() string
}

type correlation struct{}

// NewCorrelationID creates a new instance of the correlationID interface
func NewCorrelationID() Correlation {
	return correlation{}
}

// NewCorrelation will return a new correlation ID as string
func (correlation) New() string {
	return utils.NewUUID()
}
