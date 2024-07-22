package resiliency

import (
	"context"
	"fmt"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker/v2"
)

func WithCriticalError(isCritical bool, next func(context.Context) secondary.SecondaryPortCallResult) func(context.Context) secondary.SecondaryPortCallResult {
	return func(ctx context.Context) secondary.SecondaryPortCallResult {
		result := next(ctx)
		if isCritical {
			result.SetIsCritical(true)
		}
		return result
	}
}

func WithCallTimeout(timeout time.Duration, next func(context.Context) secondary.SecondaryPortCallResult) func(context.Context) secondary.SecondaryPortCallResult {
	return func(ctx context.Context) secondary.SecondaryPortCallResult {
		callCtx, callCancel := context.WithTimeout(ctx, timeout)
		result := next(callCtx)
		callCancel()
		return result
	}
}

func WithTaskTimeout(timeout time.Duration, next func(context.Context) secondary.SecondaryPortCallResult) func(context.Context) secondary.SecondaryPortCallResult {
	return func(ctx context.Context) secondary.SecondaryPortCallResult {
		callCtx, taskCancel := context.WithTimeout(ctx, timeout)
		result := next(callCtx)
		taskCancel()
		return result
	}
}

func WithRetry(spec model.RetrySpec, next func(context.Context) secondary.SecondaryPortCallResult) func(context.Context) secondary.SecondaryPortCallResult {
	return func(ctx context.Context) secondary.SecondaryPortCallResult {
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

			result = next(ctx)

			if result.Error == nil {
				return result
			}
		}
		return secondary.SecondaryPortCallResult{Payload: nil, Error: fmt.Errorf("max retry attempt exceeded, lastError=%s", result.Error)}
	}
}

func WithLogger(primaryAdapterId string, secondaryAdapter secondary.SecodaryPort, next func(context.Context) secondary.SecondaryPortCallResult) func(context.Context) secondary.SecondaryPortCallResult {
	return func(ctx context.Context) secondary.SecondaryPortCallResult {
		result := next(ctx)
		if result.Error != nil {
			logger.Logger.WithContext(ctx).Error(
				"Call failed. ",
				"sourceId=", primaryAdapterId, ", ",
				"destId=", secondaryAdapter.GetDestId(), ", ",
				"err=", result.Error,
			)
		} else {
			logger.Logger.WithContext(ctx).Debug(
				"Call succeeded. ",
				"sourceId=", primaryAdapterId, ", ",
				"destId=", secondaryAdapter.GetDestId(), ", ",
			)
		}
		return result
	}
}

func WithCircuitBreaker(spec model.CircuitBreakerSpec, secondaryAdapter secondary.SecodaryPort, next func(context.Context) secondary.SecondaryPortCallResult) func(context.Context) secondary.SecondaryPortCallResult {
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

	return func(ctx context.Context) secondary.SecondaryPortCallResult {
		result, err := cb.Execute(func() (secondary.SecondaryPortCallResult, error) {
			result := next(ctx)
			return result, result.Error
		})

		if err != nil {
			result.SetError(err)
		}

		return result
	}
}
