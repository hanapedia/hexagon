package v1

import (
	"testing"
	"time"
)

func TestGetInitialBackoff(t *testing.T) {
	tests := []struct {
		name              string
		retryInitialBackoff string
		expected          time.Duration
	}{
		{"valid duration", "100ms", 100 * time.Millisecond},
		{"invalid duration", "invalid", DEFAULT_RETRY_INITIAL_BACKOFF},
		{"empty duration", "", DEFAULT_RETRY_INITIAL_BACKOFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oes := &OnErrorSpec{
				RetryInitialBackoff: tt.retryInitialBackoff,
			}
			actual := oes.GetInitialBackoff()
			if actual != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestGetNthBackoff(t *testing.T) {
	tests := []struct {
		name               string
		policy             RetryBackoffPolicy
		initialBackoff     string
		n                  int
		expected           time.Duration
	}{
		{"constant backoff", CONSTANT_BACKOFF, "100ms", 3, 100 * time.Millisecond},
		{"linear backoff", LINEAR_BACKOFF, "100ms", 3, 300 * time.Millisecond},
		{"exponential backoff", EXPONENTIAL_BACKOFF, "100ms", 3, 400 * time.Millisecond},
		{"no backoff", NO_BACKOFF, "100ms", 3, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oes := &OnErrorSpec{
				RetryBackoffPolicy: tt.policy,
				RetryInitialBackoff: tt.initialBackoff,
			}
			actual := oes.GetNthBackoff(tt.n)
			if actual != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

