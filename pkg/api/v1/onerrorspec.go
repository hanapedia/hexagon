package v1

import (
	"github.com/hanapedia/hexagon/pkg/api/helper"
)

// OnErrorSpec is the configuration for what the service unit will do
// when the secondary adapter call fails
type OnErrorSpec struct {
	// IsCritical takes boolean specifying whether to fail
	// the parent primary adapter call when this secondary adapter call fails.
	IsCritical bool `json:"concurrent,omitempty"`

	// Retry configurations
	// RetryBackoffPolicy configurs the backoff policy used for subsequent retries.
	// Available configurations are
	//   - none: Default. No retries.
	//   - constant: Retries without any backoffs.
	//   - linear: Retries with linearly scaling backoff. Scaling can be configured from RetryBackoffScaling
	//   - exponential: Retries with exponentially scaling backoff.
	RetryPolicy string `json:"retryPolicy,omitempty"`

	// RetryMaxAttempt specifies the max number of retries before giving up.
	//
	// If both RetryMaxDuration is also specified, the one that expires first will take precedence.
	RetryMaxAttempt int `json:"retryMaxAttempt,omitempty"`

	// RetryMaxDuration specifies the max duration to retry for before giving up.
	//
	// If both RetryMaxAttempt is also specified, the one that expires first will take precedence.
	RetryMaxDuration helper.Duration `json:"retryMaxDuration,omitempty"`

	// RetryInitialBackoff specifies the initial duration to backoff
	RetryInitialBackoff helper.Duration `json:"retryInitialBackoff,omitempty"`

	// RetryBackoffScaling is the scaling factor for linear backoff policies.
	//
	// For linear backoff it is used as the multiplier for each backoff duration.
	// e.g. the backoff duration for nth attmept is
	//
	// initial backoff + backoff scaling * (n - 1) * initial backoff
	RetryBackoffScaling int             `json:"retryBackoffScaling,omitempty"`
}
