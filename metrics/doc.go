package metrics

// MetricMeta is a string for specifying the counter names
type MetricMeta struct {
	Name string
	Help string
}

// RequestCount is "request_count"
var RequestCount = MetricMeta{
	Name: "request_count",
	Help: "Number of requests received",
}

// LatencyInMili is "request_latency_microseconds"
var LatencyInMili = MetricMeta{
	Name: "request_latency_microseconds",
	Help: "Total duration of requests in microseconds",
}
