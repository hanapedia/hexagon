# Hexagon
A highly configurable microservices benchmark application generator.
Hexagon can configure each microservice in the benchmark application to expose APIs and invoke APIs exposed by other services. 

## Objective
Hexagon is a microservices benchmark suite with configurablity in mind.
It aims to solve the limited benchmarking possibilites for researchers studying microservices.
Using Hexagon, researchers can easily construct a benchmark microservices application with highly configured compositions.

## Terminology
see [terminologies](./docs/terminology.md).

## Features
### Application generation
Hexagon allows configuration of following features for generated application
#### Stateless services
Hexagon can [configure](./docs/configuration.md) APIs that each service expose and invoke using yaml.
The types of APIs include:

Synchronous protocols such as:
- HTTP REST
- gRPC

Asynchrounous protocols such as:
- Kafka

Database connections such as:
- MongoDB
- Redis

#### Service Internal Workload
Hexagon supports the configuration of internal workload for each service.
Using this feature, Hexagon can simulate resources intensive tasks performed by each service.
Currently CPU intesive workload is supported.


### Deployment manifest generation
Using the [cli](./cmd/hexctl/) Hexagon can generate kubernetes manifests files from the configuration of each service.
The generated manifests include resources such as:
- Deployment
- Service
- ConfigMap

Deployment can also configured with number of replicas for the service, resource limits and requests, and extra environmental variables.

## Components
- *service unit*: service unit is the primary container image that can be used by the stateless services.
- *hexctl*: hexctl is the cli program that can validate the configuration files and generate kubernetes deployment manifests for the configured application.

## How it works
see [internals](./docs/internals.md).

