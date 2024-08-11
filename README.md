# Hexagon
A highly configurable microservices benchmark application generator designed specifically for researchers.
Hexagon can generate benchmark microservices applications of virtually unlimited scale and interaction patterns without requiring application rebuilds.

![Hexagon Diagram](./docs/assets/hexagon_main_diagram.png)

## Quick Start
Follow [this guide](./docs/quickstart.md) to test out Hexagon, which will walk you through basic tasks for intallation and manifest generation.

## Features
The main feature of Hexagon is generating microservices application with a fine-grain control over each service.

The behaviors and interactions between services are configured by, for each service, defining a list of **primary adapters** (APIs of the service that are exposed to other services) and lists of **secondary adapters** (APIs of other services that the service interacts with) for each primary adapter.

When one of primary adapters on a service is accessed, the list of secondary adapters for that primary adapter will be invoked. Those secondary adapters will access primary adapters on other serivces, which consequentially triggers the list of secondary adapters on that service.

This would give you a full control over how services will interact with each other (and also the behavior of each API of the services).

### Different types of APIs for Adapters
Following types of APIs are currently supported for defining the primary and secondary adapters of a service.

External APIs: (Both primary and secondary adapter)
- Synchronous communication.
    - HTTP REST
    - gRPC
- Asynchrounous communication using brokers.
    - Apache Kafka
- Database connections.
    - MongoDB
    - Redis

Internal API: (Only for secondary adapter)
- Resource stressor.
    - CPU

There are no limitations on defining two or more adapters of same type. Each adapters are implemented at the endpoint granularity (e.g. HTTP route), so defining the same type of adapter multiple times on a single service just creates multiple endpoints of that type.

See [configurations](./docs/api/v1/configuration.md) for more details on how to configure each adapter.

Adding support for more APIs should be relatively easy by design, so feel free to open an issue or pull request!

### Deployment manifest generation
Using the [cli](./cmd/hexctl/) Hexagon generates Kubernetes manifests files from the configuration of each service.

These manifests can also be configured with number of replicas for the service, resource limits and requests, and extra environmental variables. See [configurations](./docs/api/v1/configuration.md) for more details.

### Resiliency patterns
Resiliency patterns such as Request Retry, Request Timeout, and Client Circuit Breaker can be configured for **each** secondary adapter. For each resiliency pattern, you can tune parameters such as retry backoffs and circuit breaker thresholds.
See [configurations](./docs/api/v1/configuration.md) for more details.

## Components
- *service unit*: service unit is the primary container image that can be used by the stateless services.
- *hexctl*: [hexctl](./docs/hexctl.md) is the cli program that can validate the configuration files and generate kubernetes deployment manifests for the configured application.

## How it works
see [internals](./docs/internals.md).

## Terminology
see [terminology](./docs/terminology.md).
