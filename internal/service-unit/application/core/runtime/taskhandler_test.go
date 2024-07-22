package runtime_test

import (
	"context"
	"testing"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/core/runtime"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/resiliency"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/test/mock"
	"github.com/sony/gobreaker/v2"
	"github.com/stretchr/testify/assert"
)

// TestRegularCalls asserts that Call methods are called properly
func TestRegularCalls(t *testing.T) {
	// 1. Prepare primary handler with tasks
	name := "RegularCallHandler"
	handler := domain.PrimaryAdapterHandler{
		ServiceName: name,
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.TaskHandler{
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					Retry:         model.RetrySpec{Disable: true},
					CircutBreaker: model.CircuitBreakerSpec{Disable: true},
				}},
				mock.NewSecondaryAdapter("RegularSecondaryAdapter1", time.Millisecond, 0),
			),
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					Retry:         model.RetrySpec{Disable: true},
					CircutBreaker: model.CircuitBreakerSpec{Disable: true},
				}},
				mock.NewSecondaryAdapter("RegularSecondaryAdapter2", time.Millisecond, 0),
			),
		},
	}

	result := runtime.TaskSetHandler(context.Background(), &handler)
	assert.False(t, result.ShouldFail)
	assert.Equal(t, 2, len(result.TaskResults))
	assert.Nil(t, result.TaskResults[0].Error)
	assert.Nil(t, result.TaskResults[1].Error)
}

// TestConcurrentCalls
func TestConcurrentCalls(t *testing.T) {
	// 1. Prepare primary handler with tasks
	name := "ConcurretCallHandler"
	handler := domain.PrimaryAdapterHandler{
		ServiceName: name,
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.TaskHandler{
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					IsCritical:    true,
					Retry:         model.RetrySpec{Disable: true},
					CircutBreaker: model.CircuitBreakerSpec{Disable: true},
				}},
				mock.NewSecondaryAdapter("ConcurrentSecondaryAdapter1", time.Millisecond, 0),
			),
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					IsCritical:    true,
					Retry:         model.RetrySpec{Disable: true},
					CircutBreaker: model.CircuitBreakerSpec{Disable: true},
				}},
				mock.NewSecondaryAdapter("ConcurrentSecondaryAdapter2", time.Millisecond, 0),
			),
		},
	}

	result := runtime.TaskSetHandler(context.Background(), &handler)
	assert.False(t, result.ShouldFail)
	assert.Equal(t, 2, len(result.TaskResults))
	assert.Nil(t, result.TaskResults[0].Error)
	assert.Nil(t, result.TaskResults[1].Error)
}

// TestRetrySuccess asserts that the call succeeds after retry
func TestRetrySuccess(t *testing.T) {
	// 1. Prepare primary handler with tasks
	name := "RetrySuccessCallHandler"
	handler := domain.PrimaryAdapterHandler{
		ServiceName: name,
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.TaskHandler{
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{
					Resiliency: model.ResiliencySpec{
						Retry: model.RetrySpec{
							BackoffPolicy:  model.NO_BACKOFF,
							MaxAttempt:     2,
							InitialBackoff: "1s",
						},
						CircutBreaker: model.CircuitBreakerSpec{Disable: true},
					},
				},
				mock.NewSecondaryAdapter("RetrySuccessSecondaryAdapter1", time.Millisecond, 1),
			),
		},
	}

	result := runtime.TaskSetHandler(context.Background(), &handler)
	assert.False(t, result.ShouldFail)
	assert.Equal(t, 1, len(result.TaskResults))
	assert.Nil(t, result.TaskResults[0].Error)
}

// TestRetryFail asserts that the call fails after retry
func TestRetryFail(t *testing.T) {
	// 1. Prepare primary handler with tasks
	name := "RetryFailCallHandler"
	handler := domain.PrimaryAdapterHandler{
		ServiceName: name,
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.TaskHandler{
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{
					Resiliency: model.ResiliencySpec{
						Retry: model.RetrySpec{
							BackoffPolicy:  model.NO_BACKOFF,
							MaxAttempt:     2,
							InitialBackoff: "1s",
						},
						CircutBreaker: model.CircuitBreakerSpec{Disable: true},
					},
				},
				mock.NewSecondaryAdapter("RetryFailSecondaryAdapter1", time.Millisecond, 5),
			),
		},
	}

	result := runtime.TaskSetHandler(context.Background(), &handler)
	assert.True(t, result.ShouldFail)
	assert.Equal(t, 1, len(result.TaskResults))
	assert.NotNil(t, result.TaskResults[0].Error)
}

// TestNonCriticalFailure
func TestNonCriticalFailure(t *testing.T) {
	// 1. Prepare primary handler with tasks
	name := "NonCriticalFailCallHandler"
	handler := domain.PrimaryAdapterHandler{
		ServiceName: name,
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.TaskHandler{
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					Retry:         model.RetrySpec{Disable: true},
					CircutBreaker: model.CircuitBreakerSpec{Disable: true},
				}},
				mock.NewSecondaryAdapter("NonCriticalSecondaryAdapter1", time.Millisecond, 1),
			),
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					Retry:         model.RetrySpec{Disable: true},
					CircutBreaker: model.CircuitBreakerSpec{Disable: true},
				}},
				mock.NewSecondaryAdapter("NonCriticalSecondaryAdapter2", time.Millisecond, 0),
			),
		},
	}

	result := runtime.TaskSetHandler(context.Background(), &handler)
	assert.False(t, result.ShouldFail)
	assert.Equal(t, 2, len(result.TaskResults))
	assert.NotNil(t, result.TaskResults[0].Error)
	assert.Nil(t, result.TaskResults[1].Error)
}

// TestCriticalFailure
func TestCriticalFailure(t *testing.T) {
	// 1. Prepare primary handler with tasks
	name := "CriticalFailCallHandler"
	handler := domain.PrimaryAdapterHandler{
		ServiceName: name,
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.TaskHandler{
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					IsCritical:    true,
					Retry:         model.RetrySpec{Disable: true},
					CircutBreaker: model.CircuitBreakerSpec{Disable: true},
				}},
				mock.NewSecondaryAdapter("CriticalSecondaryAdapter1", time.Millisecond, 1),
			),
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					Retry:         model.RetrySpec{Disable: true},
					CircutBreaker: model.CircuitBreakerSpec{Disable: true},
				}},
				mock.NewSecondaryAdapter("NonCriticalSecondaryAdapter2", time.Millisecond, 0),
			),
		},
	}

	result := runtime.TaskSetHandler(context.Background(), &handler)
	assert.True(t, result.ShouldFail)
	assert.Equal(t, 2, len(result.TaskResults))
	assert.NotNil(t, result.TaskResults[0].Error)
	assert.Nil(t, result.TaskResults[1].Error)
}

// TestCallTimeoutFailure
func TestCallTimeoutFailure(t *testing.T) {
	// 1. Prepare primary handler with tasks
	name := "CallTimeoutFailCallHandler"
	handler := domain.PrimaryAdapterHandler{
		ServiceName: name,
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.TaskHandler{
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					CallTimeout:   "5ms",
					Retry:         model.RetrySpec{Disable: true},
					CircutBreaker: model.CircuitBreakerSpec{Disable: true},
				}},
				mock.NewSecondaryAdapter("TimeoutFailSecondaryAdapter1", 15*time.Millisecond, 0),
			),
		},
	}

	result := runtime.TaskSetHandler(context.Background(), &handler)
	assert.True(t, result.ShouldFail)
	assert.Equal(t, 1, len(result.TaskResults))
	assert.NotNil(t, result.TaskResults[0].Error)
}

// TestTaskTimeoutFailure
func TestTaskTimeoutFailure(t *testing.T) {
	// 1. Prepare primary handler with tasks
	name := "TaskTimeoutFailCallHandler"
	handler := domain.PrimaryAdapterHandler{
		ServiceName: name,
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.TaskHandler{
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					TaskTimeout: "1s",
					Retry: model.RetrySpec{
						BackoffPolicy:  model.NO_BACKOFF,
						MaxAttempt:     3,
						InitialBackoff: "1s",
					},
					CircutBreaker: model.CircuitBreakerSpec{Disable: true},
				}},
				mock.NewSecondaryAdapter("TimeoutFailSecondaryAdapter1", time.Millisecond, 5),
			),
		},
	}

	result := runtime.TaskSetHandler(context.Background(), &handler)
	assert.True(t, result.ShouldFail)
	assert.Equal(t, 1, len(result.TaskResults))
	assert.NotNil(t, result.TaskResults[0].Error)
}

func TestCircuitBreakerWithoutThresh(t *testing.T) {
	// 1. Prepare primary handler with tasks
	name := "CircutBreakerWithoutThresh"
	adapter := mock.NewSecondaryAdapter("CircuitBreakerSecondaryAdapter1", time.Millisecond, 1)
	adapter.SetDestId("mock.dest.1")
	ratio := 0
	minRequests := 0
	consecutiveFails := 0
	handler := domain.PrimaryAdapterHandler{
		ServiceName: name,
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.TaskHandler{
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					CircutBreaker: model.CircuitBreakerSpec{
						MaxRequests:      1,
						Interval:         "10s",
						Timeout:          "50ms",
						MinRequests:      uint32(minRequests),
						Ratio:            float64(ratio),
						ConsecutiveFails: uint32(consecutiveFails),
					},
					Retry: model.RetrySpec{Disable: true},
				}},
				adapter,
			),
		},
	}

	// Fail once to open circuit breaker
	_ = runtime.TaskSetHandler(context.Background(), &handler)

	for i := 0; i < 3; i++ {
		result := runtime.TaskSetHandler(context.Background(), &handler)
		assert.True(t, result.ShouldFail)
		assert.Equal(t, gobreaker.ErrOpenState, result.TaskResults[0].Error)
	}

	// Wait till circuit breaker half opens
	time.Sleep(100 * time.Millisecond)

	result2 := runtime.TaskSetHandler(context.Background(), &handler)

	assert.False(t, result2.ShouldFail)
	assert.Nil(t, result2.TaskResults[0].Error)
}

func TestCircuitBreakerRatioThresh(t *testing.T) {
	// 1. Prepare primary handler with tasks
	name := "CircutBreakerRatioThresh"
	// first three calls fail
	adapter := mock.NewSecondaryAdapter("CircuitBreakerSecondaryAdapter1", time.Millisecond, 3)
	adapter.SetDestId("mock.dest.1")
	minRequests := 3
	ratio := 0.9
	consecutiveFails := 0
	handler := domain.PrimaryAdapterHandler{
		ServiceName: name,
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.TaskHandler{
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					CircutBreaker: model.CircuitBreakerSpec{
						MaxRequests:      1,
						Interval:         "10s",
						Timeout:          "50ms",
						MinRequests:      uint32(minRequests),
						Ratio:            float64(ratio),
						ConsecutiveFails: uint32(consecutiveFails),
					},
					Retry: model.RetrySpec{Disable: true},
				}},
				adapter,
			),
		},
	}

	// Fail once to open circuit breaker
	// First three calls fail
	for i := 0; i < 3; i++ {
		_ = runtime.TaskSetHandler(context.Background(), &handler)
	}
	// Circuit Breaker should be open now
	result1 := runtime.TaskSetHandler(context.Background(), &handler)
	assert.True(t, result1.ShouldFail)
	assert.Equal(t, gobreaker.ErrOpenState, result1.TaskResults[0].Error)

	// Wait till circuit breaker half opens
	time.Sleep(100 * time.Millisecond)

	result2 := runtime.TaskSetHandler(context.Background(), &handler)

	assert.False(t, result2.ShouldFail)
	assert.Nil(t, result2.TaskResults[0].Error)
}

func TestCircuitBreakerConsecutiveFailsThresh(t *testing.T) {
	// 1. Prepare primary handler with tasks
	name := "CircuitBreakerConsecutiveFailsThresh"
	// first three calls fail
	adapter := mock.NewSecondaryAdapter("CircuitBreakerSecondaryAdapter1", time.Millisecond, 3)
	adapter.SetDestId("mock.dest.1")
	minRequests := 3
	ratio := 0
	consecutiveFails := 3
	handler := domain.PrimaryAdapterHandler{
		ServiceName: name,
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.TaskHandler{
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					CircutBreaker: model.CircuitBreakerSpec{
						MaxRequests:      1,
						Interval:         "10s",
						Timeout:          "50ms",
						MinRequests:      uint32(minRequests),
						Ratio:            float64(ratio),
						ConsecutiveFails: uint32(consecutiveFails),
					},
					Retry: model.RetrySpec{Disable: true},
				}},
				adapter,
			),
		},
	}

	// Fail once to open circuit breaker
	// First three calls fail
	for i := 0; i < 3; i++ {
		_ = runtime.TaskSetHandler(context.Background(), &handler)
	}
	// Circuit Breaker should be open now
	result1 := runtime.TaskSetHandler(context.Background(), &handler)
	assert.True(t, result1.ShouldFail)
	assert.Equal(t, gobreaker.ErrOpenState, result1.TaskResults[0].Error)

	// Wait till circuit breaker half opens
	time.Sleep(100 * time.Millisecond)

	result2 := runtime.TaskSetHandler(context.Background(), &handler)

	assert.False(t, result2.ShouldFail)
	assert.Nil(t, result2.TaskResults[0].Error)
}

func TestCircuitBreakerHalfOpen(t *testing.T) {
	// 1. Prepare primary handler with tasks
	name := "CircuitBreakerHalfOpen"
	// first three calls fail
	adapter := mock.NewSecondaryAdapter("CircuitBreakerSecondaryAdapter1", time.Millisecond, 4)
	adapter.SetDestId("mock.dest.1")
	minRequests := 3
	ratio := 0
	consecutiveFails := 3
	handler := domain.PrimaryAdapterHandler{
		ServiceName: name,
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.TaskHandler{
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					CircutBreaker: model.CircuitBreakerSpec{
						MaxRequests:      1,
						Interval:         "10s",
						Timeout:          "50ms",
						MinRequests:      uint32(minRequests),
						Ratio:            float64(ratio),
						ConsecutiveFails: uint32(consecutiveFails),
					},
					Retry: model.RetrySpec{Disable: true},
				}},
				adapter,
			),
		},
	}

	// Fail once to open circuit breaker
	// First three calls fail
	for i := 0; i < 3; i++ {
		_ = runtime.TaskSetHandler(context.Background(), &handler)
	}
	// Circuit Breaker should be open now
	result1 := runtime.TaskSetHandler(context.Background(), &handler)
	assert.True(t, result1.ShouldFail)
	assert.Equal(t, gobreaker.ErrOpenState, result1.TaskResults[0].Error)

	// Wait till circuit breaker half opens
	time.Sleep(100 * time.Millisecond)

	// Should re-open circuit breaker
	result2 := runtime.TaskSetHandler(context.Background(), &handler)
	assert.True(t, result2.ShouldFail)
	assert.NotNil(t, result2.TaskResults[0].Error)

	result3 := runtime.TaskSetHandler(context.Background(), &handler)
	assert.True(t, result3.ShouldFail)
	assert.Equal(t, gobreaker.ErrOpenState, result3.TaskResults[0].Error)

	// Wait till circuit breaker half opens
	time.Sleep(100 * time.Millisecond)

	result4 := runtime.TaskSetHandler(context.Background(), &handler)
	assert.False(t, result4.ShouldFail)
	assert.Nil(t, result4.TaskResults[0].Error)
}

func TestRetryThenCircuitBreak(t *testing.T) {
	// 1. Prepare primary handler with tasks
	name := "RetryThenCircuitBreak"
	// first three calls fail
	adapter := mock.NewSecondaryAdapter("RetryCircuitBreakerSecondaryAdapter1", time.Millisecond, 5)
	adapter.SetDestId("mock.dest.1")
	minRequests := 0
	ratio := 0
	consecutiveFails := 0
	countRetries := false
	handler := domain.PrimaryAdapterHandler{
		ServiceName: name,
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.TaskHandler{
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					CircutBreaker: model.CircuitBreakerSpec{
						MaxRequests:      1,
						Interval:         "10s",
						Timeout:          "50ms",
						MinRequests:      uint32(minRequests),
						Ratio:            float64(ratio),
						ConsecutiveFails: uint32(consecutiveFails),
						CountRetries:     countRetries,
					},
					Retry: model.RetrySpec{
						BackoffPolicy:  model.NO_BACKOFF,
						MaxAttempt:     3,
						InitialBackoff: "1ms",
					},
				}},
				adapter,
			),
		},
	}

	// Fail once to open circuit breaker
	// retries 3 times + 1st attmpt, each attempt should fail not because of circuit breaker
	// circuit breaker should open after retries are done and fail
	_ = runtime.TaskSetHandler(context.Background(), &handler)

	// Wait till circuit breaker half opens
	time.Sleep(100 * time.Millisecond)

	// Circuit Breaker should be half open now
	// retries once and succeeds
	// should not re-open
	result1 := runtime.TaskSetHandler(context.Background(), &handler)
	assert.False(t, result1.ShouldFail)
	assert.Nil(t, result1.TaskResults[0].Error)
}

func TestCircuitBreakThenRetry(t *testing.T) {
	// 1. Prepare primary handler with tasks
	name := "CircuitBreakThenRetry"
	// first three calls fail
	adapter := mock.NewSecondaryAdapter("CircuitBreakerRetrySecondaryAdapter1", time.Millisecond, 5)
	adapter.SetDestId("mock.dest.1")
	minRequests := 0
	ratio := 0
	consecutiveFails := 0
	countRetries := true
	handler := domain.PrimaryAdapterHandler{
		ServiceName: name,
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.TaskHandler{
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{
					CircutBreaker: model.CircuitBreakerSpec{
						MaxRequests:      1,
						Interval:         "10s",
						Timeout:          "50ms",
						MinRequests:      uint32(minRequests),
						Ratio:            float64(ratio),
						ConsecutiveFails: uint32(consecutiveFails),
						CountRetries:     countRetries,
					},
					Retry: model.RetrySpec{
						BackoffPolicy:  model.NO_BACKOFF,
						MaxAttempt:     3,
						InitialBackoff: "1ms",
					},
				}},
				adapter,
			),
		},
	}

	// Fail once to open circuit breaker
	// 1st attempt should open circuit breaker
	// retries 3 times, they should all fail due to circuit breaker
	// circuit breaker should open after retries are done and fail
	_ = runtime.TaskSetHandler(context.Background(), &handler)

	// Wait till circuit breaker half opens
	time.Sleep(100 * time.Millisecond)

	// Circuit Breaker should be half open now
	// retries once and fail
	// should re-open
	result1 := runtime.TaskSetHandler(context.Background(), &handler)
	assert.True(t, result1.ShouldFail)
	assert.NotNil(t, result1.TaskResults[0].Error)
}
