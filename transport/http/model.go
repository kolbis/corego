package http

// Response is a base response which will be returned from the transport
type Response struct {
	// Error which occured on the callee
	Error error

	// Data payload returned from the callee
	Data interface{}

	// CircuitOpened is true if execution failed due to circuit being opened
	IsCircuitOpen bool

	// RetryFailed is true if retry execution failed
	IsReachedRateLimit bool

	// StatusCode for the execution
	StatusCode int
}

// Request is a base response which will be send to the transport
type Request struct {
	// Specific request data
	Data interface{} `mapstructure:",squash"`

	// the correlation ID
	CorrelationID string
}
