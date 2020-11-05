package loadbalancer

import "time"

const (
	// DefaultMaxAttempts defines number of retry attempts per request, before giving up
	DefaultMaxAttempts = 5

	// DefaultRetryTimeout is the default time we are willing to wait and retry
	DefaultRetryTimeout time.Duration = time.Second * 15
)
