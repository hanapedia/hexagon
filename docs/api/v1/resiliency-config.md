## Resiliency Config
These configurations can only be applied to each secondary adapter individually.
The definitions for Go struct is in [`pkg/api/v1/resiliencyspec.go`](../../../pkg/api/v1/resiliencyspec.go)

- call timeout: timeout for calling the secondary adapter **once**.
- task timeout: timeout for calling the secondary adapter **possibly multiple times**. Timeout for all the attempts.

| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| isCritical     | When set to true, error on this secondary adapter call will result in primary adapter error | false | false |
| retry         | configurations for [retry](#retry) | {} | false |
| circuitBreaker| configurations for [circuit breaker](#circuit-breaker) | {} | false |
| callTimeout   | configurations for call timeout. Must be parsable by `time.ParseDuration` | "" | false |
| taskTimeout   | configurations for task timeout. Must be parsable by `time.ParseDuration` | "" | false |

## Retry
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| disabled      | Disable retry | false | false |
| backoffPolicy | Backoff scaling policy for retry intervals. Can be `none`, `constant`, `linear`, `exponential` | 'none' | false |
| maxAttempt    | Maximum number of attempts before giving up | 0 | false |
| initialBackoff| Initial backoff interval. Subsequent intervals will scale off of this value. Must be parsable by `time.ParseDuration`. | 1ms | false |

## Circuit Breaker
- `maxRequests`, `interval`, `timeout` are passed straight to the [gobreaker circuit breaker](https://github.com/sony/gobreaker).
- if both `ratio` and `consecutiveFails` are set, both are checked with `or` clause, meaning that condition is met when either of them is met.

| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| disabled      | Disable circuit breaker | false | false |
| maxRequests   | Maximum number of requests to allow go through when circuit breaker is half-open | 0 (set to 1 by gobreaker internally) | false |
| interval      | Interval duration to clear internal counts for threshold, Must be parsable by `time.ParseDuration` | "" (gobreaker does not clear internal counts) | false |
| timeout       | Timeout duration to wait in open state before attempting to close circuit breaker, Must be parsable by `time.ParseDuration` | "" (defaults to 60 seconds) | false |
| minRequests   | Minimum number of requests required for checking threshold | 0 | false |
| ratio         | Threshold ratio of failed request / total number of requests | 0 | false |
| consecutiveFails| The number of consecutive failures as threshold | 0 | false |
| countRetries | whether to count retries towards the circuit breaker threshold | false | false |
