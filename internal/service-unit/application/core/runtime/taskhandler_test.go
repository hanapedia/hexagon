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
				model.TaskSpec{},
				mock.NewSecondaryAdapter("RegularSecondaryAdapter1", time.Second, 0),
			),
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{},
				mock.NewSecondaryAdapter("RegularSecondaryAdapter2", time.Second, 0),
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
				model.TaskSpec{Resiliency: model.ResiliencySpec{IsCritical: true}},
				mock.NewSecondaryAdapter("ConcurrentSecondaryAdapter1", time.Second, 0),
			),
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{Resiliency: model.ResiliencySpec{IsCritical: true}},
				mock.NewSecondaryAdapter("ConcurrentSecondaryAdapter2", time.Second, 0),
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
					},
				},
				mock.NewSecondaryAdapter("RetrySuccessSecondaryAdapter1", time.Second, 1),
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
					},
				},
				mock.NewSecondaryAdapter("RetryFailSecondaryAdapter1", time.Second, 5),
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
				model.TaskSpec{},
				mock.NewSecondaryAdapter("NonCriticalSecondaryAdapter1", time.Second, 1),
			),
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{},
				mock.NewSecondaryAdapter("NonCriticalSecondaryAdapter2", time.Second, 0),
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
				model.TaskSpec{
					Resiliency: model.ResiliencySpec{IsCritical: true},
				},
				mock.NewSecondaryAdapter("CriticalSecondaryAdapter1", time.Second, 1),
			),
			resiliency.NewTaskHandler(
				name,
				model.TaskSpec{},
				mock.NewSecondaryAdapter("NonCriticalSecondaryAdapter2", time.Second, 0),
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
				model.TaskSpec{
					Resiliency: model.ResiliencySpec{
						CallTimeout: "1s",
					},
				},
				mock.NewSecondaryAdapter("TimeoutFailSecondaryAdapter1", 2*time.Second, 0),
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
				model.TaskSpec{
					Resiliency: model.ResiliencySpec{
						TaskTimeout: "2s",
						Retry: model.RetrySpec{
							BackoffPolicy:  model.NO_BACKOFF,
							MaxAttempt:     3,
							InitialBackoff: "1s",
						},
					},
				},
				mock.NewSecondaryAdapter("TimeoutFailSecondaryAdapter1", time.Second, 5),
			),
		},
	}

	result := runtime.TaskSetHandler(context.Background(), &handler)
	assert.True(t, result.ShouldFail)
	assert.Equal(t, 1, len(result.TaskResults))
	assert.NotNil(t, result.TaskResults[0].Error)
}
