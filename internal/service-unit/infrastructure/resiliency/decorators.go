package resiliency

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hanapedia/adapto/rto"
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
type CallWithContextAlias = func(context.Context, *TaskContext) secondary.SecondaryPortCallResult
type CircuitBreaker = *gobreaker.CircuitBreaker[secondary.SecondaryPortCallResult]

type TaskContext struct {
	mu             sync.Mutex
	circuitBreaker CircuitBreaker
	attempt        uint32
	telemetryCtx   domain.TelemetryContext
	isConcurrent   bool
}

func (tc *TaskContext) IncAttempt() {
	tc.mu.Lock()
	tc.attempt += 1
	tc.mu.Unlock()
}

// WithUnWrapTaskContext decorates `next` function with a decorator that unwraps TaskContext into context.Context
// Should be the first decorator
func WithUnWrapTaskContext(next CallAlias) CallWithContextAlias {
	return func(ctx context.Context, taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		return next(ctx)
	}
}

func WithCallDurationMetrics(next CallWithContextAlias, async bool) CallWithContextAlias {
	return func(ctx context.Context, tc *TaskContext) secondary.SecondaryPortCallResult {
		startTime := time.Now()
		result := next(ctx, tc)
		elapsed := time.Since(startTime)

		var status domain.Status
		switch result.Error {
		case nil:
			status = domain.Ok
		case context.DeadlineExceeded:
			status = domain.ErrCtxDeadlineExceeded
		case context.Canceled:
			status = domain.ErrCtxCanceled
		case gobreaker.ErrOpenState:
			status = domain.ErrCBOpen
		case rto.RequestRateLimitExceeded:
			status = domain.ErrCBOpen
		default:
			status = domain.ErrGeneric
		}

		var circuitBreakerState = ""
		if tc.circuitBreaker != nil {
			circuitBreakerState = tc.circuitBreaker.State().String()
		}

		if async {
			go metrics.ObserveSecondaryAdapterCallDuration(elapsed, domain.SecondaryAdapterCallDurationLabels{
				Ctx:                 tc.telemetryCtx,
				Status:              status,
				NthAttmpt:           tc.attempt,
				CircuitBreakerState: circuitBreakerState,
				IsConcurrent:        tc.isConcurrent,
			})
		} else {
			metrics.ObserveSecondaryAdapterCallDuration(elapsed, domain.SecondaryAdapterCallDurationLabels{
				Ctx:                 tc.telemetryCtx,
				Status:              status,
				NthAttmpt:           tc.attempt,
				CircuitBreakerState: circuitBreakerState,
				IsConcurrent:        tc.isConcurrent,
			})
		}
		return result
	}
}

func WithTaskDurationMetrics(next CallWithContextAlias, async bool) CallWithContextAlias {
	return func(ctx context.Context, tc *TaskContext) secondary.SecondaryPortCallResult {
		startTime := time.Now()
		result := next(ctx, tc)
		elapsed := time.Since(startTime)

		var status domain.Status
		switch result.Error {
		case nil:
			status = domain.Ok
		case context.DeadlineExceeded:
			status = domain.ErrCtxDeadlineExceeded
		case context.Canceled:
			status = domain.ErrCtxCanceled
		case gobreaker.ErrOpenState:
			status = domain.ErrCBOpen
		case rto.RequestRateLimitExceeded:
			status = domain.ErrCBOpen
		default:
			status = domain.ErrGeneric
		}

		var circuitBreakerState = ""
		if tc.circuitBreaker != nil {
			circuitBreakerState = tc.circuitBreaker.State().String()
		}

		if async {
			go metrics.ObserveSecondaryAdapterTaskDuration(elapsed, domain.SecondaryAdapterTaskDurationLabels{
				Ctx:                 tc.telemetryCtx,
				Status:              status,
				TotalAttempts:       tc.attempt,
				CircuitBreakerState: circuitBreakerState,
				IsConcurrent:        tc.isConcurrent,
			})
		} else {
			metrics.ObserveSecondaryAdapterTaskDuration(elapsed, domain.SecondaryAdapterTaskDurationLabels{
				Ctx:                 tc.telemetryCtx,
				Status:              status,
				TotalAttempts:       tc.attempt,
				CircuitBreakerState: circuitBreakerState,
				IsConcurrent:        tc.isConcurrent,
			})
		}
		return result
	}
}

func WithCriticalError(isCritical bool, next CallWithContextAlias) CallWithContextAlias {
	return func(ctx context.Context, taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		result := next(ctx, taskCtx)
		if isCritical {
			result.SetIsCritical(true)
		}
		return result
	}
}

func WithCallTimeout(timeout time.Duration, next CallWithContextAlias) CallWithContextAlias {
	return func(ctx context.Context, taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		newCtx, callCancel := context.WithTimeout(ctx, timeout)
		defer callCancel()
		result := next(newCtx, taskCtx)
		return result
	}
}

func WithTaskTimeout(timeout time.Duration, next CallWithContextAlias) CallWithContextAlias {
	return func(ctx context.Context, taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		newCtx, taskCancel := context.WithTimeout(ctx, timeout)
		defer taskCancel()
		result := next(newCtx, taskCtx)
		return result
	}
}

func WithAdaptiveRTOCallTimeout(spec model.AdaptiveTimeoutSpec, secondaryAdapter secondary.SecodaryPort, next CallWithContextAlias) CallWithContextAlias {
	adaptoRTOConfig := rto.Config{
		Id:             secondaryAdapter.GetDestId(),
		Min:            spec.GetMin(),
		Max:            spec.GetMax(),
		SLOFailureRate: spec.SLO,
		/* Capacity:       spec.Capacity, */
		Interval: spec.GetInterval(),
		KMargin:  spec.KMargin,
		Logger:   logger.AdaptoLogger,
	}
	return func(ctx context.Context, taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		timeoutDuration, rttCh, err := rto.GetTimeout(ctx, adaptoRTOConfig)
		if err != nil {
			if err == rto.RequestRateLimitExceeded {
				return secondary.SecondaryPortCallResult{
					Payload: nil,
					Error:   err,
				}
			}
			logger.Logger.
				WithField("id", adaptoRTOConfig.Id).
				WithField("err", err).
				Errorf("failed to create new adaptive RTO timeout. resorting to default call timeout.")
			timeoutDuration = model.DEFAULT_CALL_TIMEOUT
		}
		// record gauge metrics for timeout value
		metrics.SetAdaptiveCallTimeoutDuration(timeoutDuration, domain.AdaptiveTimeoutGaugeLabels{Ctx: taskCtx.telemetryCtx})

		newCtx, callCancel := context.WithTimeout(ctx, timeoutDuration)
		defer callCancel()

		startTime := time.Now()
		result := next(newCtx, taskCtx)

		// Check for cancelation or timeout after `next` returns
		if newCtx.Err() == context.DeadlineExceeded {
			// send negative timeout duration to signal that deadline exceeded
			// from adapto v1.0.14
			rttCh <- -timeoutDuration
		} else if result.Error == nil {
			elapsed := time.Since(startTime)
			rttCh <- elapsed
		}
		return result
	}
}

func WithRetry(spec model.RetrySpec, next CallWithContextAlias) CallWithContextAlias {
	return func(ctx context.Context, taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		// add 1 for the initial attempt
		maxAttempt := spec.MaxAttempt + 1
		var result secondary.SecondaryPortCallResult

		for i := 0; i < maxAttempt; i++ {
			if i > 0 {
				backoff := spec.GetNthBackoff(i)
				timer := time.NewTimer(backoff)
				select {
				// check for the parent context expiration
				case <-ctx.Done():
					timer.Stop()
					return secondary.SecondaryPortCallResult{Payload: nil, Error: ctx.Err()}
				case <-timer.C:
					timer.Stop()
				}
			}

			taskCtx.IncAttempt()
			result = next(ctx, taskCtx)

			if result.Error == nil {
				return result
			}
		}
		return secondary.SecondaryPortCallResult{Payload: nil, Error: fmt.Errorf("max retry attempt exceeded, lastError=%s", result.Error)}
	}
}

func WithLogger(telCtx domain.TelemetryContext, timing string, secondaryAdapter secondary.SecodaryPort, next CallWithContextAlias) CallWithContextAlias {
	return func(ctx context.Context, taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		result := next(ctx, taskCtx)
		if result.Error != nil {
			logger.Logger.
				WithField("timing", timing).
				WithField("destId", secondaryAdapter.GetDestId()).
				WithField("err", result.Error).
				Error("Call failed.")
		} else {
			logger.Logger.
				WithField("timing", timing).
				WithField("destId", secondaryAdapter.GetDestId()).
				Debug("Call succeeded.")
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

	return func(ctx context.Context, taskCtx *TaskContext) secondary.SecondaryPortCallResult {
		result, err := cb.Execute(func() (secondary.SecondaryPortCallResult, error) {
			result := next(ctx, taskCtx)
			return result, result.Error
		})

		if err != nil {
			result.SetError(err)
		}

		return result
	}, cb
}
