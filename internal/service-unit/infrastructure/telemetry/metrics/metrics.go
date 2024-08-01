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
		prometheus.MustRegister(instance.PrimaryAdapterDuration)
		prometheus.MustRegister(instance.SecondaryAdapterCallDuration)
		prometheus.MustRegister(instance.SecondaryAdapterTaskDuration)
		prometheus.MustRegister(instance.CallTimeout)
		prometheus.MustRegister(instance.TaskTimeout)
		prometheus.MustRegister(instance.CircuitBreakerDisabled)
		prometheus.MustRegister(instance.CircuitBreakerCountRetries)
		prometheus.MustRegister(instance.CircuitBreakerIntervalSecs)
		prometheus.MustRegister(instance.CircuitBreakerMaxRequests)
		prometheus.MustRegister(instance.CircuitBreakerMinRequests)
		prometheus.MustRegister(instance.CircuitBreakerTimeout)
		prometheus.MustRegister(instance.CircuitBreakerRatio)
		prometheus.MustRegister(instance.CircuitBreakerConsecutiveFails)
		prometheus.MustRegister(instance.RetryDisabled)
		prometheus.MustRegister(instance.RetryMaxAttempt)
		prometheus.MustRegister(instance.RetryInitialBackoff)
	})
	return instance
}
