package tracer

import (
	"os"

	zipkingo "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

// Tracer ...
type Tracer struct {
	ServiceName       string
	ServerHostAddress string
	URL               string

	Inst *zipkingo.Tracer
}

// NewTracer ...
func NewTracer(serviceName string, hostAddress string, zipkinURL string) Tracer {
	return Tracer{
		ServiceName:       serviceName,
		ServerHostAddress: hostAddress,
		URL:               zipkinURL,
	}
}

// Instance will return the zipkin single instance
func (t *Tracer) Instance() *zipkingo.Tracer {
	if t.Inst == nil {
		var zipkinTracer *zipkingo.Tracer
		{
			if t.URL != "" {
				var (
					err         error
					hostPort    = t.ServerHostAddress
					serviceName = t.ServiceName
					reporter    = zipkinhttp.NewReporter(t.URL)
				)
				defer reporter.Close()
				zEP, _ := zipkingo.NewEndpoint(serviceName, hostPort)
				zipkinTracer, err = zipkingo.NewTracer(reporter, zipkingo.WithLocalEndpoint(zEP))

				if err != nil {
					os.Exit(1)
				}
			}
		}

		t.Inst = zipkinTracer
	}

	return t.Inst
}
