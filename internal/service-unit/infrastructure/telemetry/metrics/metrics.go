package metrics

import (
	"sync"

	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	PrimaryAdapterDuration         *prometheus.HistogramVec
	SecondaryAdapterCallDuration   *prometheus.HistogramVec
	SecondaryAdapterTaskDuration   *prometheus.HistogramVec
	PrimaryAdapterInProgress       *prometheus.GaugeVec
	AdaptiveTaskTimeoutDuration    *prometheus.GaugeVec
	AdaptiveCallTimeoutDuration    *prometheus.GaugeVec
	CallTimeout                    *prometheus.GaugeVec
	TaskTimeout                    *prometheus.GaugeVec
	CircuitBreakerDisabled         *prometheus.GaugeVec
	CircuitBreakerCountRetries     *prometheus.GaugeVec
	CircuitBreakerIntervalSecs     *prometheus.GaugeVec
	CircuitBreakerMaxRequests      *prometheus.GaugeVec
	CircuitBreakerMinRequests      *prometheus.GaugeVec
	CircuitBreakerTimeout          *prometheus.GaugeVec
	CircuitBreakerRatio            *prometheus.GaugeVec
	CircuitBreakerConsecutiveFails *prometheus.GaugeVec
	RetryDisabled                  *prometheus.GaugeVec
	RetryMaxAttempt                *prometheus.GaugeVec
	RetryInitialBackoff            *prometheus.GaugeVec
}

var (
	instance *Metrics
	once     sync.Once
)

func GetInstance() *Metrics {
	once.Do(func() {
		instance = &Metrics{
			PrimaryAdapterDuration:         domain.GetHistogramVec(domain.PrimaryAdapterDuration),
			SecondaryAdapterCallDuration:   domain.GetHistogramVec(domain.SecondaryAdapterCallDuration),
			SecondaryAdapterTaskDuration:   domain.GetHistogramVec(domain.SecondaryAdapterTaskDuration),
			PrimaryAdapterInProgress:       domain.GetGaugeVec(domain.PrimaryAdapterInProgress),
			AdaptiveTaskTimeoutDuration:    domain.GetGaugeVec(domain.AdaptiveTaskTimeoutDuration),
			AdaptiveCallTimeoutDuration:    domain.GetGaugeVec(domain.AdaptiveCallTimeoutDuration),
			CallTimeout:                    domain.GetGaugeVec(domain.CallTimeout),
			TaskTimeout:                    domain.GetGaugeVec(domain.TaskTimeout),
			CircuitBreakerDisabled:         domain.GetGaugeVec(domain.CircuitBreakerDisabled),
			CircuitBreakerCountRetries:     domain.GetGaugeVec(domain.CircuitBreakerCountRetries),
			CircuitBreakerIntervalSecs:     domain.GetGaugeVec(domain.CircuitBreakerIntervalSecs),
			CircuitBreakerMaxRequests:      domain.GetGaugeVec(domain.CircuitBreakerMaxRequests),
			CircuitBreakerMinRequests:      domain.GetGaugeVec(domain.CircuitBreakerMinRequests),
			CircuitBreakerTimeout:          domain.GetGaugeVec(domain.CircuitBreakerTimeout),
			CircuitBreakerRatio:            domain.GetGaugeVec(domain.CircuitBreakerRatio),
			CircuitBreakerConsecutiveFails: domain.GetGaugeVec(domain.CircuitBreakerConsecutiveFails),
			RetryDisabled:                  domain.GetGaugeVec(domain.RetryDisabled),
			RetryMaxAttempt:                domain.GetGaugeVec(domain.RetryMaxAttempt),
			RetryInitialBackoff:            domain.GetGaugeVec(domain.RetryInitialBackoff),
		}

		// Register the metrics
		PromRegistry.MustRegister(instance.PrimaryAdapterDuration)
		PromRegistry.MustRegister(instance.SecondaryAdapterCallDuration)
		PromRegistry.MustRegister(instance.SecondaryAdapterTaskDuration)
		PromRegistry.MustRegister(instance.PrimaryAdapterInProgress)
		PromRegistry.MustRegister(instance.AdaptiveTaskTimeoutDuration)
		PromRegistry.MustRegister(instance.AdaptiveCallTimeoutDuration)
		PromRegistry.MustRegister(instance.CallTimeout)
		PromRegistry.MustRegister(instance.TaskTimeout)
		PromRegistry.MustRegister(instance.CircuitBreakerDisabled)
		PromRegistry.MustRegister(instance.CircuitBreakerCountRetries)
		PromRegistry.MustRegister(instance.CircuitBreakerIntervalSecs)
		PromRegistry.MustRegister(instance.CircuitBreakerMaxRequests)
		PromRegistry.MustRegister(instance.CircuitBreakerMinRequests)
		PromRegistry.MustRegister(instance.CircuitBreakerTimeout)
		PromRegistry.MustRegister(instance.CircuitBreakerRatio)
		PromRegistry.MustRegister(instance.CircuitBreakerConsecutiveFails)
		PromRegistry.MustRegister(instance.RetryDisabled)
		PromRegistry.MustRegister(instance.RetryMaxAttempt)
		PromRegistry.MustRegister(instance.RetryInitialBackoff)
	})
	return instance
}
