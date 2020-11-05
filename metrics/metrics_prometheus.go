package metrics

import (
	metrics "github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// PrometheusInstrumentor ..
type PrometheusInstrumentor struct {
	ServiceName        string
	PromCounters       []metrics.Counter
	PromSummaryVectors []metrics.Histogram
	PromGauges         []metrics.Histogram
	PromHistograms     []metrics.Histogram
}

// NewPrometheusInstrumentor ...
func NewPrometheusInstrumentor(serviceName string) PrometheusInstrumentor {
	return PrometheusInstrumentor{
		ServiceName: serviceName,
	}
}

// AddPromCounter will add a new counter to prometheous
// Namespace, Subsystem, and Name are components of the fully-qualified name of the Metric (created by joining these components with "_").
// Only Name is mandatory, the others merely help structuring the name.
// promLabels are labels to differentiate the characteristics of the thing that is being measured
// read here for more information: https://prometheus.io/docs/practices/naming/
func (inst *PrometheusInstrumentor) AddPromCounter(namespace string, subsystem string, info MetricMeta, promLabels []string) metrics.Counter {
	opts := stdprometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      info.Name,
		Help:      info.Help,
	}
	var counter metrics.Counter
	counter = kitprometheus.NewCounterFrom(opts, promLabels)
	inst.PromCounters = append(inst.PromCounters, counter)
	return counter
}

// AddPromSummary ...
func (inst *PrometheusInstrumentor) AddPromSummary(namespace string, subsystem string, info MetricMeta, promLabels []string) metrics.Histogram {
	opts := stdprometheus.SummaryOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      info.Name,
		Help:      info.Help,
	}

	var summary metrics.Histogram
	summary = kitprometheus.NewSummaryFrom(opts, promLabels)
	inst.PromSummaryVectors = append(inst.PromSummaryVectors, summary)
	return summary
}

// AddPromGauge ...
func (inst *PrometheusInstrumentor) AddPromGauge(namespace string, subsystem string, info MetricMeta, promLabels []string) {
	// TBD
}

// AddPromHistogram ...
func (inst *PrometheusInstrumentor) AddPromHistogram(namespace string, subsystem string, info MetricMeta, promLabels []string) {
	// TBD
}
