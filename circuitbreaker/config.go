package circuitbreaker

var (
	// DefaultTimeout is how long to wait for command to complete, in milliseconds
	DefaultTimeout = 1000
	// DefaultMaxConcurrent is how many commands of the same type can run at the same time
	DefaultMaxConcurrent = 10
	// DefaultVolumeThreshold is the minimum number of requests needed before a circuit can be tripped due to health
	DefaultVolumeThreshold = 20
	// DefaultSleepWindow is how long, in milliseconds, to wait after a circuit opens before testing for recovery
	DefaultSleepWindow = 5000
	// DefaultErrorPercentThreshold causes circuits to open once the rolling measure of errors exceeds this percent of requests
	DefaultErrorPercentThreshold = 50
)

// HystrixCommandConfig ...
type HystrixCommandConfig struct {
	Timeout                int `json:"timeout"`
	MaxConcurrentRequests  int `json:"max_concurrent_requests"`
	RequestVolumeThreshold int `json:"request_volume_threshold"`
	SleepWindow            int `json:"sleep_window"`
	ErrorPercentThreshold  int `json:"error_percent_threshold"`
}

// NewHystrixCommandConfig will return a new HystrixCommandConfig with default values
func NewHystrixCommandConfig() HystrixCommandConfig {
	return HystrixCommandConfig{
		MaxConcurrentRequests:  DefaultMaxConcurrent,
		ErrorPercentThreshold:  DefaultErrorPercentThreshold,
		Timeout:                DefaultTimeout,
		RequestVolumeThreshold: DefaultVolumeThreshold,
		SleepWindow:            DefaultSleepWindow,
	}
}
