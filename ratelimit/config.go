package ratelimit

import (
	"time"

	"golang.org/x/time/rate"
)

// RateLimiterConfig ...
type RateLimiterConfig struct {
	MaxQueriesPerSecond int
	Limit               rate.Limit
}

const (
	defaultMaxQueriesPerSecond int = 100
)

var (
	defaultLimit rate.Limit = rate.Every(time.Second)
)

// NewDefaultRateLimiterConfig ...
func NewDefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		MaxQueriesPerSecond: defaultMaxQueriesPerSecond,
		Limit:               defaultLimit,
	}
}
