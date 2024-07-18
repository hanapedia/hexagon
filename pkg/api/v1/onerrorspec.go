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

// OnErrorSpec is the configuration for what the service unit will do
// when the secondary adapter call fails
type OnErrorSpec struct {
	// IsCritical takes boolean specifying whether to fail
	// the parent primary adapter call when this secondary adapter call fails.
	IsCritical bool `json:"concurrent,omitempty"`

	// Retry configurations
	Retry RetrySpec `json:"retry,omitempty"`
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
	//
	// If both RetryMaxDuration is also specified, the one that expires first will take precedence.
	MaxAttempt int `json:"maxAttempt,omitempty"`

	// InitialBackoff specifies the initial duration to backoff
	//
	// Must be parsable using time.ParseDuration
	InitialBackoff string `json:"initialBackoff,omitempty"`
}

// Get parsed timeout as time.Duration
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
