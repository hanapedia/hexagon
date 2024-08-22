// Groundtruth for available metrics and label names.
package domain

import (
	"fmt"
	"slices"

	"github.com/hanapedia/hexagon/pkg/operator/utils"
	"github.com/prometheus/client_golang/prometheus"
)

type GaugeOp int

const (
	INC GaugeOp = iota
	DEC
)

type GaugeVecName = string
type HistogramVecName = string

const (
	PrimaryAdapterDuration         HistogramVecName = "primary_adapter_duration_ms"
	SecondaryAdapterCallDuration   HistogramVecName = "secondary_adapter_call_duration_ms"
	SecondaryAdapterTaskDuration   HistogramVecName = "secondary_adapter_task_duration_ms"
	PrimaryAdapterInProgress       GaugeVecName     = "primary_adapter_in_progress"
	CallTimeout                    GaugeVecName     = "call_timeout_ms"
	TaskTimeout                    GaugeVecName     = "task_timeout_ms"
	CircuitBreakerDisabled         GaugeVecName     = "circuit_breaker_disabled"
	CircuitBreakerCountRetries     GaugeVecName     = "circuit_breaker_count_retries"
	CircuitBreakerIntervalSecs     GaugeVecName     = "circuit_breaker_interval_seconds"
	CircuitBreakerMaxRequests      GaugeVecName     = "circuit_breaker_max_requests"
	CircuitBreakerMinRequests      GaugeVecName     = "circuit_breaker_min_requests"
	CircuitBreakerTimeout          GaugeVecName     = "circuit_breaker_timeout_seconds"
	CircuitBreakerRatio            GaugeVecName     = "circuit_breaker_ratio"
	CircuitBreakerConsecutiveFails GaugeVecName     = "circuit_breaker_consecutive_fails"
	RetryDisabled                  GaugeVecName     = "retry_disabled"
	RetryMaxAttempt                GaugeVecName     = "retry_max_attempt"
	RetryInitialBackoff            GaugeVecName     = "retry_initial_backoff"
)

// Default Histogram bucket is []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
// and desinged as response time in seconds
// Adjust it so that it is in milliseconds
var DurationHistogramBuckets = slices.Concat(
	[]float64{0.1, 0.25, 0.5},
	[]float64{1, 2.5, 5},
	[]float64{10, 25, 50},
	[]float64{100, 250, 500},
	[]float64{1000, 2500, 5000, 10000},
)

var histogramVecs map[HistogramVecName]*prometheus.HistogramVec = map[string]*prometheus.HistogramVec{
	PrimaryAdapterDuration: prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    PrimaryAdapterDuration,
			Help:    "Duration for primary adapter execution in milliseconds.",
			Buckets: DurationHistogramBuckets,
		},
		utils.GetMapKeys(PrimaryAdapterDurationLabels{}.AsMap()),
	),
	SecondaryAdapterCallDuration: prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    SecondaryAdapterCallDuration,
			Help:    "Duration for secondary adapter call execution in milliseconds.",
			Buckets: DurationHistogramBuckets,
		},
		utils.GetMapKeys(SecondaryAdapterCallDurationLabels{}.AsMap()),
	),
	SecondaryAdapterTaskDuration: prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    SecondaryAdapterTaskDuration,
			Help:    "Duration for secondary adapter task execution in milliseconds.",
			Buckets: DurationHistogramBuckets,
		},
		utils.GetMapKeys(SecondaryAdapterTaskDurationLabels{}.AsMap()),
	),
}

var gaugeVecs map[GaugeVecName]*prometheus.GaugeVec = map[GaugeVecName]*prometheus.GaugeVec{
	PrimaryAdapterInProgress: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: PrimaryAdapterInProgress,
			Help: "Gauge for primary adapter in progress requests.",
		},
		utils.GetMapKeys(PrimaryAdapterInProgressLabels{}.AsMap()),
	),
	CallTimeout: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: CallTimeout,
			Help: "Timeout for each secondary adapter call.",
		},
		utils.GetMapKeys(TimeoutGaugeLabels{}.AsMap()),
	),
	TaskTimeout: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: TaskTimeout,
			Help: "Timeout for each secondary adapter task.",
		},
		utils.GetMapKeys(TimeoutGaugeLabels{}.AsMap()),
	),
	CircuitBreakerDisabled: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: CircuitBreakerDisabled,
			Help: "Flag that indicates whether circuit breaker is disabled.",
		},
		utils.GetMapKeys(CircuitBreakerGaugeLabels{}.AsMap()),
	),
	CircuitBreakerCountRetries: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: CircuitBreakerCountRetries,
			Help: "Flag that indicates whether circuit breaker counts retries.",
		},
		utils.GetMapKeys(CircuitBreakerGaugeLabels{}.AsMap()),
	),
	CircuitBreakerIntervalSecs: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: CircuitBreakerIntervalSecs,
			Help: "Interval duration for circuit breaker cycles in seconds.",
		},
		utils.GetMapKeys(CircuitBreakerGaugeLabels{}.AsMap()),
	),
	CircuitBreakerMaxRequests: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: CircuitBreakerMaxRequests,
			Help: "Max number of requests that circuit breaker allows when in half-open state.",
		},
		utils.GetMapKeys(CircuitBreakerGaugeLabels{}.AsMap()),
	),
	CircuitBreakerMinRequests: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: CircuitBreakerMinRequests,
			Help: "Min number of requests during interval that circuit breaker needs to trip.",
		},
		utils.GetMapKeys(CircuitBreakerGaugeLabels{}.AsMap()),
	),
	CircuitBreakerTimeout: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: CircuitBreakerTimeout,
			Help: "Timeout for circuit breaker to wait in open state in seconds.",
		},
		utils.GetMapKeys(CircuitBreakerGaugeLabels{}.AsMap()),
	),
	CircuitBreakerRatio: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: CircuitBreakerRatio,
			Help: "Threshold for circuit breaker given as ratio of failed requests / total requests.",
		},
		utils.GetMapKeys(CircuitBreakerGaugeLabels{}.AsMap()),
	),
	CircuitBreakerConsecutiveFails: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: CircuitBreakerConsecutiveFails,
			Help: "Threshold for circuit breaker given as number of consecutive fails.",
		},
		utils.GetMapKeys(CircuitBreakerGaugeLabels{}.AsMap()),
	),
	RetryDisabled: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: RetryDisabled,
			Help: "Flag that indicates whether retry is disabled.",
		},
		utils.GetMapKeys(RetryGaugeLabels{}.AsMap()),
	),
	RetryMaxAttempt: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: RetryMaxAttempt,
			Help: "Max retry attempts allowed.",
		},
		utils.GetMapKeys(RetryGaugeLabels{}.AsMap()),
	),
	RetryInitialBackoff: prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: RetryInitialBackoff,
			Help: "Initial backoff duration for retries.",
		},
		utils.GetMapKeys(RetryGaugeLabels{}.AsMap()),
	),
}

func GetHistogramVec(name HistogramVecName) *prometheus.HistogramVec {
	histogramVec, ok := histogramVecs[name]
	if !ok {
		panic(fmt.Sprintf("No such histogram metrics availbale. name=%v", name))
	}
	return histogramVec
}

func GetGaugeVec(name GaugeVecName) *prometheus.GaugeVec {
	gaugeVec, ok := gaugeVecs[name]
	if !ok {
		panic(fmt.Sprintf("No such gauge metrics availbale. name=%v", name))
	}
	return gaugeVec
}

type PrimaryAdapterDurationLabels struct {
	PrimaryLabels
	Status Status
}

func (padl PrimaryAdapterDurationLabels) AsMap() map[string]string {
	base := padl.GetPrimaryLabels()
	base["status"] = padl.Status.AsString()
	return base
}

type PrimaryAdapterInProgressLabels struct {
	PrimaryLabels
}

func (p PrimaryAdapterInProgressLabels) AsMap() map[string]string {
	base := p.GetPrimaryLabels()
	return base
}

type SecondaryAdapterCallDurationLabels struct {
	Ctx                 TelemetryContext
	Status              Status
	NthAttmpt           uint32
	CircuitBreakerState string
	IsConcurrent        bool
}

func (sacdl SecondaryAdapterCallDurationLabels) AsMap() map[string]string {
	base := sacdl.Ctx.AsMap()
	base["status"] = sacdl.Status.AsString()
	base["nth_attempt"] = fmt.Sprint(sacdl.NthAttmpt)
	base["circuit_breaker_state"] = sacdl.CircuitBreakerState
	base["is_concurrent"] = utils.Btos(sacdl.IsConcurrent)
	return base
}

type SecondaryAdapterTaskDurationLabels struct {
	Ctx                 TelemetryContext
	Status              Status
	TotalAttempts       uint32
	CircuitBreakerState string
	IsConcurrent        bool
}

func (satdl SecondaryAdapterTaskDurationLabels) AsMap() map[string]string {
	base := satdl.Ctx.AsMap()
	base["status"] = satdl.Status.AsString()
	base["total_attempts"] = fmt.Sprint(satdl.TotalAttempts)
	base["circuit_breaker_state"] = satdl.CircuitBreakerState
	base["is_concurrent"] = utils.Btos(satdl.IsConcurrent)
	return base
}

type TimeoutGaugeLabels struct {
	Ctx TelemetryContext
}

func (tgl TimeoutGaugeLabels) AsMap() map[string]string {
	return tgl.Ctx.AsMap()
}

type CircuitBreakerGaugeLabels struct {
	Ctx TelemetryContext
}

func (cbgl CircuitBreakerGaugeLabels) AsMap() map[string]string {
	return cbgl.Ctx.AsMap()
}

type RetryGaugeLabels struct {
	Ctx           TelemetryContext
	BackoffPolicy string
}

func (rgl RetryGaugeLabels) AsMap() map[string]string {
	base := rgl.Ctx.AsMap()
	base["backoff_policy"] = rgl.BackoffPolicy
	return base
}
