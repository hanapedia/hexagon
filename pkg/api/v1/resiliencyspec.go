package v1

import (
	"math"
	"time"
)

type RetryBackoffPolicy = string

const (
	NO_BACKOFF          RetryBackoffPolicy = "none"
	CONSTANT_BACKOFF    RetryBackoffPolicy = "constant"
	LINEAR_BACKOFF      RetryBackoffPolicy = "linear"
	EXPONENTIAL_BACKOFF RetryBackoffPolicy = "exponential"
)

var (
	DEFAULT_RETRY_INITIAL_BACKOFF time.Duration = time.Millisecond
)

var (
	DEFAULT_CIRCUIT_BREAKER_INTERVAL     time.Duration = time.Minute
	DEFAULT_CIRCUIT_BREAKER_TIMEOUT      time.Duration = time.Minute
)

var (
	DEFAULT_TASK_TIMEOUT = 5 * time.Second
	DEFAULT_CALL_TIMEOUT = 2 * time.Second
)

// ResiliencySpec is the configuration for what the service unit will do
// when the secondary adapter call fails
type ResiliencySpec struct {
	// IsCritical takes boolean specifying whether to fail
	// the parent primary adapter call when this secondary adapter call fails.
	IsCritical bool `json:"concurrent,omitempty"`

	// Retry configurations
	Retry RetrySpec `json:"retry,omitempty"`

	// CircutBreaker configurations
	CircutBreaker CircuitBreakerSpec `json:"circuitBreaker,omitempty"`

	// taskTimeout is used as the value for request timeout of the calls INCLUDING all retries.
	// Must be parsable with time.ParseDuration, otherwise default value will be used.
	TaskTimeout string `json:"taskTimeout,omitempty"`
	// callTimeout refers to the timeout assigend to each call attempt
	CallTimeout string `json:"callTimeout,omitempty"`
}

type RetrySpec struct {
	// BackoffPolicy configurs the backoff policy used for subsequent retries.
	// Available configurations are
	//   - none: Default. No retry backoff.
	//   - constant: Retries without any backoffs.
	//   - linear: Retries with linearly scaling backoff. Scaling can be configured from RetryBackoffScaling
	//   - exponential: Retries with exponentially scaling backoff.
	BackoffPolicy RetryBackoffPolicy `json:"backoffPolicy,omitempty"`

	// MaxAttempt specifies the max number of retries before giving up.
	MaxAttempt int `json:"maxAttempt,omitempty"`

	// InitialBackoff specifies the initial duration to backoff
	//
	// Must be parsable using time.ParseDuration
	InitialBackoff string `json:"initialBackoff,omitempty"`

	// Disable is the flag to turn off retry feature entirely
	Disable bool `json:"disable,omitempty"`
}

type CircuitBreakerSpec struct {
	// MaxRequests is the max number of requests to go through when half open
	// if MaxRequests is 0, circuit breaker will only allow 1 request
	// Use GetMaxRequests to read the value with default values
	MaxRequests uint32 `json:"maxRequests,omitempty"`

	// Interval is the cycle duration to clear the internal count for request success / failure during closed state
	// Must be parsable with time.ParseDuration, otherwise default value will be used
	Interval string `json:"interval,omitempty"`

	// Timeout is the amount of duration that the circuit breaker stays open after breaching the threshold
	// Must be parsable with time.ParseDuration, otherwise default value will be used
	Timeout string `json:"timeout,omitempty"`

	// MinRequests is the min number of requests required to open the circuit
	MinRequests uint32 `json:"minRequests,omitempty"`

	// Ratio is the threshold ratio of failed request / total number of requests
	// If ConsecutiveFails is set as well, whichever trips first will take precedence
	// If neither is set, circuit breaker will trip on the first error
	Ratio float64 `json:"ratio,omitempty"`

	// ConsecutiveFails is the number of consecutive failures as threshold
	// If Ratio is set as well, whichever trips first will take precedence
	// If neither is set, circuit breaker will trip on the first error
	ConsecutiveFails uint32 `json:"consecutiveFails,omitempty"`

	// CountRetries specify whether the retry attempts are counted by the circuit breaker.
	// if set to `true`, the circuit breaker counts each attempt.
	// if set to `false`, the circuit breaker counts all retry attempts as a single request.
	// Note that circuit broken requests will also be retried when set to true.
	CountRetries bool `json:"countRetries,omitempty"`

	// Disable is the flag to turn off circuit breaker feature entirely
	Disable bool `json:"disable,omitempty"`
}

// Get parsed initial backoff as time.Duration
func (rs *RetrySpec) GetInitialBackoff() time.Duration {
	duration, err := time.ParseDuration(rs.InitialBackoff)
	if err != nil {
		return DEFAULT_RETRY_INITIAL_BACKOFF
	}
	return duration
}

// GetNthBackoff returns the backoff duration for Nth retry attempt.
func (rs *RetrySpec) GetNthBackoff(n int) time.Duration {
	switch rs.BackoffPolicy {
	case CONSTANT_BACKOFF:
		return rs.GetInitialBackoff()
	case LINEAR_BACKOFF:
		return rs.GetInitialBackoff() * time.Duration(n)
	case EXPONENTIAL_BACKOFF:
		return rs.GetInitialBackoff() * time.Duration(math.Pow(2, float64(n-1)))
	case NO_BACKOFF:
	}
	return 0
}

// Get parsed circuit breaker interval as time.Duration
func (cbs *CircuitBreakerSpec) GetInterval() time.Duration {
	duration, err := time.ParseDuration(cbs.Interval)
	if err != nil {
		return DEFAULT_CIRCUIT_BREAKER_INTERVAL
	}
	return duration
}

// Get parsed circuit breaker timeout as time.Duration
func (cbs *CircuitBreakerSpec) GetTimeout() time.Duration {
	duration, err := time.ParseDuration(cbs.Timeout)
	if err != nil {
		return DEFAULT_CIRCUIT_BREAKER_TIMEOUT
	}
	return duration
}

// Get parsed taskTimeout as time.Duration
func (r ResiliencySpec) GetTaskTimeout() time.Duration {
	duration, err := time.ParseDuration(r.TaskTimeout)
	if err != nil {
		return DEFAULT_TASK_TIMEOUT
	}
	if duration == 0 {
		return DEFAULT_TASK_TIMEOUT
	}
	return duration
}

// Get parsed getCallTimeout as time.Duration
func (r ResiliencySpec) GetCallTimeout() time.Duration {
	duration, err := time.ParseDuration(r.CallTimeout)
	if err != nil {
		return DEFAULT_CALL_TIMEOUT
	}
	if duration == 0 {
		return DEFAULT_CALL_TIMEOUT
	}
	return duration
}
