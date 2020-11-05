package logger

type nopLogger struct{}

// NewNopLogger returns a log.Logger that doesn't do anything.
// Should be used `for testing only
func NewNopLogger() Logger {
	return &nopLogger{}
}

func (logger nopLogger) Log(keyvals ...interface{}) error {
	return nil
}
