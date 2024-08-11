## Contents
This directory contains documentations for all the available configuration values.

**Some content may or may not be up to date, but most should be.**

see [terminologies](./docs/terminology.md) for unclear terms.

## Basics
- All the Configurations should be provided in YAML.

### [Service Unit Configurations](./service-unit-config.md)
Configuration values for each Service Unit.
- primary adapters
- secondary adapters

### [Deployment Configurations](./deployment-config.md)
Deployment configurations for each Service Unit.
- replicas
- resources
- environment variables

### [Resiliency Configurations](./resiliency-config.md)
Resiliency Configurations for each secondary adapter of Service Unit.
- retry
- timeout
- circuit breaker

### [Cluster Configurations](./cluster-config.md)
Global configuration for all Service Unit.
- namespace
- port numbers
