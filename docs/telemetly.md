## Telemetry
Metrics and Distributed Tracing Spans can be collected from each Service Unit.

### Metrics
The metrics are measured natively by hexagon and exported via metrics server running at `:7070/metrics` (default).

Each metrics are collected with labels that you can used to identify the adapter.
#### Histogram
Histograms for duration metrics.
- `primary_adapter_duration_ms`
- `secondary_adapter_call_duration_ms`
- `secondary_adapter_task_duration_ms`
#### Gauge
Gauge metrics for configuration values.
- `call_timeout_ms`
- `task_timeout_ms`
- `circuit_breaker_disabled`
- `circuit_breaker_count_retries`
- `circuit_breaker_interval_seconds`
- `circuit_breaker_max_requests`
- `circuit_breaker_min_requests`
- `circuit_breaker_timeout`
- `circuit_breaker_ratio`
- `circuit_breaker_consecutive_fails`
- `retry_disabled`
- `retry_max_attempt`
- `retry_intial_backoff_ms`

### Distributed Tracing Spans
Trace spans are generated for each adapter calls and collected using otel sdk and exported via Otel Collector.
Otel Collector must be configured separately, and if not, tracing is disabled at runtime.

### Logs
Logs are currently not centralized, but can be observe vie `kube-apiserver` using `kubectl logs`
