package loadbalancer

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
)

// DynamicLoadBalancer has dynamic set of endpoints populated from service discovery
type DynamicLoadBalancer struct {
	Endpointer *sd.DefaultEndpointer
}

// NewDynamicLoadBalancer creates a new instance of DynamicLoadBalancer
func NewDynamicLoadBalancer(endpointer *sd.DefaultEndpointer) DynamicLoadBalancer {
	return DynamicLoadBalancer{
		Endpointer: endpointer,
	}
}

// DefaultRoundRobinWithRetryEndpoint ...
func (b *DynamicLoadBalancer) DefaultRoundRobinWithRetryEndpoint(ctx context.Context) endpoint.Endpoint {
	return b.RoundRobinWithRetryEndpoint(DefaultMaxAttempts, DefaultRetryTimeout)
}

// RoundRobinWithRetryEndpoint ..
func (b *DynamicLoadBalancer) RoundRobinWithRetryEndpoint(maxAttempts int, maxTime time.Duration) endpoint.Endpoint {
	balancer := lb.NewRoundRobin(b.Endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)
	return retry
}
