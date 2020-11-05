package circuitbreaker_test

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/endpoint"
	cb "github.com/kolbis/corego/circuitbreaker"
)

func TestHystrixCommandMiddleware(t *testing.T) {
	const (
		commandName   = "my-endpoint"
		errorPercent  = 5
		maxConcurrent = 1000
	)

	config := cb.HystrixCommandConfig{
		ErrorPercentThreshold: errorPercent,
		MaxConcurrentRequests: maxConcurrent,
	}

	var (
		breaker    = cb.NewHystrixCommandMiddleware(commandName, config)
		primeWith  = hystrix.DefaultVolumeThreshold * 2
		shouldPass = func(n int) bool { return (float64(n) / float64(primeWith+n)) <= (float64(errorPercent-1) / 100.0) }
	)
	// hystrix-go uses buffered channels to receive reports on request success/failure,
	// and so is basically impossible to test deterministically. We have to make sure
	// the report buffer is emptied, by injecting a sleep between each invocation.
	requestDelay := 5 * time.Millisecond

	testFailingEndpoint(t, breaker, primeWith, shouldPass, requestDelay)
}

func testFailingEndpoint(
	t *testing.T,
	breaker endpoint.Middleware,
	primeWith int,
	shouldPass func(int) bool,
	requestDelay time.Duration,
) {
	_, file, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s:%d", filepath.Base(file), line)

	// Create a mock endpoint and wrap it with the breaker.
	m := mock{}
	var e endpoint.Endpoint
	e = m.endpoint
	e = breaker(e)

	// Prime the endpoint with successful requests.
	for i := 0; i < primeWith; i++ {
		if _, err := e(context.Background(), struct{}{}); err != nil {
			t.Fatalf("%s: during priming, got error: %v", caller, err)
		}
		time.Sleep(requestDelay)
	}

	// Switch the endpoint to start throwing errors.
	m.err = errors.New("tragedy+disaster")
	m.through = 0

	// The first several should be allowed through and yield our error.
	for i := 0; shouldPass(i); i++ {
		if _, err := e(context.Background(), struct{}{}); err != m.err {
			t.Fatalf("%s: want %v, have %v", caller, m.err, err)
		}
		time.Sleep(requestDelay)
	}
	through := m.through

	// But the rest should be blocked by an open circuit.
	for i := 0; i < 10; i++ {
		if _, err := e(context.Background(), struct{}{}); err.Error() != hystrix.ErrCircuitOpen.Error() {
			t.Fatalf("%s: want %q, have %q", caller, hystrix.ErrCircuitOpen.Error(), err.Error())
		}
		time.Sleep(requestDelay)
	}

	// Make sure none of those got through.
	if want, have := through, m.through; want != have {
		t.Errorf("%s: want %d, have %d", caller, want, have)
	}
}

type mock struct {
	through int
	err     error
}

func (m *mock) endpoint(context.Context, interface{}) (interface{}, error) {
	m.through++
	return struct{}{}, m.err
}
