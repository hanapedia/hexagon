package runtime_test

import (
	"context"
	"testing"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/core/runtime"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/test/mock"
	"github.com/stretchr/testify/assert"
)

// TestRegularCalls asserts that Call methods are called properly
func TestRegularCalls(t *testing.T) {
	// 1. Prepare primary handler with tasks
	handler := domain.PrimaryHandler{
		ServiceName: "RegularCallHandler",
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.Task{
			{
				SecondaryPort: mock.NewSecondaryAdapter("RegularSecondaryAdapter1", time.Second, 0),
			},
			{
				SecondaryPort: mock.NewSecondaryAdapter("RegularSecondaryAdapter2", time.Second, 0),
			},
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
	handler := domain.PrimaryHandler{
		ServiceName: "ConcurretCallHandler",
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.Task{
			{
				SecondaryPort: mock.NewSecondaryAdapter("ConcurrentSecondaryAdapter1", time.Second, 0),
				Concurrent:    true,
			},
			{
				SecondaryPort: mock.NewSecondaryAdapter("ConcurrentSecondaryAdapter2", time.Second, 0),
				Concurrent:    true,
			},
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
	handler := domain.PrimaryHandler{
		ServiceName: "RetrySuccessCallHandler",
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.Task{
			{
				SecondaryPort: mock.NewSecondaryAdapter("RetrySuccessSecondaryAdapter1", time.Second, 1),
				OnError: model.OnErrorSpec{
					Retry: model.RetrySpec{
						BackoffPolicy:  model.NO_BACKOFF,
						MaxAttempt:     2,
						InitialBackoff: "1s",
					},
				},
			},
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
	handler := domain.PrimaryHandler{
		ServiceName: "RetryFailCallHandler",
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.Task{
			{
				SecondaryPort: mock.NewSecondaryAdapter("RetryFailSecondaryAdapter1", time.Second, 5),
				OnError: model.OnErrorSpec{
					Retry: model.RetrySpec{
						BackoffPolicy:  model.NO_BACKOFF,
						MaxAttempt:     2,
						InitialBackoff: "1s",
					},
				},
			},
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
	handler := domain.PrimaryHandler{
		ServiceName: "NonCriticalFailCallHandler",
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.Task{
			{
				SecondaryPort: mock.NewSecondaryAdapter("NonCriticalSecondaryAdapter1", time.Second, 1),
			},
			{
				SecondaryPort: mock.NewSecondaryAdapter("NonCriticalSecondaryAdapter2", time.Second, 0),
			},
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
	handler := domain.PrimaryHandler{
		ServiceName: "CriticalFailCallHandler",
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.Task{
			{
				SecondaryPort: mock.NewSecondaryAdapter("CriticalSecondaryAdapter1", time.Second, 1),
				OnError:       model.OnErrorSpec{IsCritical: true},
			},
			{
				SecondaryPort: mock.NewSecondaryAdapter("NonCriticalSecondaryAdapter1", time.Second, 0),
			},
		},
	}

	result := runtime.TaskSetHandler(context.Background(), &handler)
	assert.True(t, result.ShouldFail)
	assert.Equal(t, 2, len(result.TaskResults))
	assert.NotNil(t, result.TaskResults[0].Error)
	assert.Nil(t, result.TaskResults[1].Error)
}

// TestTimeoutFailure
func TestTimeoutFailure(t *testing.T) {
	// 1. Prepare primary handler with tasks
	handler := domain.PrimaryHandler{
		ServiceName: "TimeoutFailCallHandler",
		ServerConfig: &model.ServerConfig{
			Variant: constants.REST,
			Action:  constants.GET,
			Route:   "test",
		},
		TaskSet: []domain.Task{
			{
				SecondaryPort: mock.NewSecondaryAdapter("TimeoutFailSecondaryAdapter1", 2*time.Second, 0),
				CallTimeout:   "1s",
			},
		},
	}

	result := runtime.TaskSetHandler(context.Background(), &handler)
	assert.True(t, result.ShouldFail)
	assert.Equal(t, 1, len(result.TaskResults))
	assert.NotNil(t, result.TaskResults[0].Error)
}
