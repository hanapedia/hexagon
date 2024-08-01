package resiliency

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/metrics"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker/v2"
)

type alias = secondary.SecondaryPortCallResult
type CallAlias = func(context.Context) secondary.SecondaryPortCallResult
type CallWithContextAlias = func(*TaskContext) secondary.SecondaryPortCallResult
type CircuitBreaker = *gobreaker.CircuitBreaker[secondary.SecondaryPortCallResult]

type TaskContext struct {
	mu             sync.Mutex
	ctx            context.Context
	circuitBreaker CircuitBreaker
	attempt        uint32
	telemetryCtx   domain.TelemetryContext
	isConcurrent   bool
}

func (tc *TaskContext) WithTimeout(duration time.Duration) context.CancelFunc {
	newCtx, cancel := context.WithTimeout(tc.ctx, duration)

	tc.mu.Lock()
	tc.ctx = newCtx
	tc.mu.Unlock()

	return cancel
}

func (tc *TaskContext) IncAttempt() {
	tc.mu.Lock()
	tc.attempt += 1
	tc.mu.Unlock()
}

// WithUnWrapTaskContext decorates `next` function with a decorator that unwraps TaskContext into context.Context
// Should be the first decorator
func WithUnWrapTaskContext(next CallAlias) CallWithContextAlias {
	return func(taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		return next(taskCtx.ctx)
	}
}

func WithCallDurationMetrics(next CallWithContextAlias) CallWithContextAlias {
	return func(tc *TaskContext) secondary.SecondaryPortCallResult {
		startTime := time.Now()
		result := next(tc)
		elapsed := time.Since(startTime)

		go func() {
			var status domain.Status
			if result.Error == nil {
				status = domain.Ok
			} else {
				status = domain.Err
			}

			var circuitBreakerState = ""
			if tc.circuitBreaker != nil {
				circuitBreakerState = tc.circuitBreaker.State().String()
			}

			metrics.ObserveSecondaryAdapterCallDuration(elapsed, domain.SecondaryAdapterCallDurationLabels{
				Ctx:                 tc.telemetryCtx,
				Status:              status,
				NthAttmpt:           tc.attempt,
				CircuitBreakerState: circuitBreakerState,
				IsConcurrent:        tc.isConcurrent,
			})
		}()
		return result
	}
}

func WithTaskDurationMetrics(next CallWithContextAlias) CallWithContextAlias {
	return func(tc *TaskContext) secondary.SecondaryPortCallResult {
		startTime := time.Now()
		result := next(tc)
		elapsed := time.Since(startTime)

		go func() {
			var status domain.Status
			if result.Error == nil {
				status = domain.Ok
			} else {
				status = domain.Err
			}

			var circuitBreakerState = ""
			if tc.circuitBreaker != nil {
				circuitBreakerState = tc.circuitBreaker.State().String()
			}

			metrics.ObserveSecondaryAdapterTaskDuration(elapsed, domain.SecondaryAdapterTaskDurationLabels{
				Ctx:                 tc.telemetryCtx,
				Status:              status,
				TotalAttempts:       tc.attempt,
				CircuitBreakerState: circuitBreakerState,
				IsConcurrent:        tc.isConcurrent,
			})
		}()
		return result
	}
}

func WithCriticalError(isCritical bool, next CallWithContextAlias) CallWithContextAlias {
	return func(taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		result := next(taskCtx)
		if isCritical {
			result.SetIsCritical(true)
		}
		return result
	}
}

func WithCallTimeout(timeout time.Duration, next CallWithContextAlias) CallWithContextAlias {
	return func(taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		callCancel := taskCtx.WithTimeout(timeout)
		result := next(taskCtx)
		callCancel()
		return result
	}
}

func WithTaskTimeout(timeout time.Duration, next CallWithContextAlias) CallWithContextAlias {
	return func(taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		taskCancel := taskCtx.WithTimeout(timeout)
		result := next(taskCtx)
		taskCancel()
		return result
	}
}

func WithRetry(spec model.RetrySpec, next CallWithContextAlias) CallWithContextAlias {
	return func(taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		// add 1 for the initial attempt
		maxAttempt := spec.MaxAttempt + 1
		var result secondary.SecondaryPortCallResult

		for i := 0; i < maxAttempt; i++ {
			if i > 0 {
				backoff := spec.GetNthBackoff(i)
				timer := time.NewTimer(backoff)
				select {
				// check for the parent context expiration
				case <-taskCtx.ctx.Done():
					timer.Stop()
					return secondary.SecondaryPortCallResult{Payload: nil, Error: taskCtx.ctx.Err()}
				case <-timer.C:
					timer.Stop()
				}
			}

			taskCtx.IncAttempt()
			result = next(taskCtx)

			if result.Error == nil {
				return result
			}
		}
		return secondary.SecondaryPortCallResult{Payload: nil, Error: fmt.Errorf("max retry attempt exceeded, lastError=%s", result.Error)}
	}
}

func WithLogger(telCtx domain.TelemetryContext, secondaryAdapter secondary.SecodaryPort, next CallWithContextAlias) CallWithContextAlias {
	return func(taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		result := next(taskCtx)
		if result.Error != nil {
			logger.Logger.Error(
				"Call failed. ",
				/* "sourceId=", primaryAdapterId, ", ", */
				"destId=", secondaryAdapter.GetDestId(), ", ",
				"err=", result.Error,
			)
		} else {
			logger.Logger.Debug(
				"Call succeeded. ",
				/* "sourceId=", primaryAdapterId, ", ", */
				"destId=", secondaryAdapter.GetDestId(), ", ",
			)
		}
		return result
	}
}

func WithCircuitBreaker(spec model.CircuitBreakerSpec, secondaryAdapter secondary.SecodaryPort, next CallWithContextAlias) (CallWithContextAlias, CircuitBreaker) {
	setting := gobreaker.Settings{
		Name:        secondaryAdapter.GetDestId(),
		MaxRequests: spec.MaxRequests,
		Interval:    spec.GetInterval(),
		Timeout:     spec.GetTimeout(),
		// Ratio threshold and ConsecutiveFails threshold is considered at the same time and which ever trips first will take precedence
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			consecutiveFailsThresh := spec.ConsecutiveFails > 0 && counts.ConsecutiveFailures >= spec.ConsecutiveFails
			ratioThresh := failureRatio > spec.Ratio
			return counts.Requests >= spec.MinRequests && (ratioThresh || consecutiveFailsThresh)
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			logger.Logger.WithFields(logrus.Fields{
				"name": name,
				"from": from.String(),
				"to":   to.String(),
			}).Info("Circuit Breaker State Updated")
		},
	}
	cb := gobreaker.NewCircuitBreaker[secondary.SecondaryPortCallResult](setting)

	return func(taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		result, err := cb.Execute(func() (secondary.SecondaryPortCallResult, error) {
			result := next(taskCtx)
			return result, result.Error
		})

		if err != nil {
			result.SetError(err)
		}

		return result
	}, cb
}
