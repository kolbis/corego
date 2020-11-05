package loadbalancer

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
)

// FixedLoadBalancer has a fix set of endpoints to load balance between
type FixedLoadBalancer struct {
	Endpointer sd.FixedEndpointer
}

// NewFixedLoadBalancer ..
func NewFixedLoadBalancer(fixedEndpointer sd.FixedEndpointer) FixedLoadBalancer {
	return FixedLoadBalancer{
		Endpointer: fixedEndpointer,
	}
}

// DefaultRoundRobinWithRetryEndpoint ...
func (b *FixedLoadBalancer) DefaultRoundRobinWithRetryEndpoint(ctx context.Context) endpoint.Endpoint {
	return b.RoundRobinWithRetryEndpoint(DefaultMaxAttempts, DefaultRetryTimeout)
}

// RoundRobinWithRetryEndpoint ..
func (b *FixedLoadBalancer) RoundRobinWithRetryEndpoint(maxAttempts int, maxTime time.Duration) endpoint.Endpoint {
	balancer := lb.NewRoundRobin(b.Endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)
	return retry
}
