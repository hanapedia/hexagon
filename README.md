# the-bench
A highly configurable microservices benchmark application.
Each microservice in the benchmark application can be configured with their own sets of ingress and egress handling logic. 

## Objective
the-bench is a microservices benchmark suite with configurablity in mind.
It aims to solve the limited benchmarking possibilites for researchers studying microservices.
Using the-bench researchers can easily construct a benchmark microservices application with highly configured compositions.

## Terminology
see [terminologies](./docs/terminology.md).

## Features
### Application generation
the-bench allows configuration of following features for generated application
#### Stateless services
Each stateless service in the-bench can be [configured](./docs/configuration.md) to serve and make various types of network protocols.
The types of network protocols include:

Synchronous protocols such as:
- HTTP REST
- gRPC

Asynchrounous protocols such as:
- Kafka

Database connections such as:
- MongoDB
- PostgresQL [WIP]
- Redis [WIP]

#### Stateful Service
the-bench supports the creation and configuration of stateful services. Currently supported services are:
- MongoDB
- Kakfa

#### Service Internal Workload [WIP]
the-bench supports the configuration of internal workload for each service.
Using this feature, any stateless services can simulate various types of resources intensive tasks.


### Deployment manifest generation
Using the [cli](./cmd/tbctl/) the-bench can generate kubernetes manifests files from the configuration of each service.
The generated manifests include resources such as:
- Deployment
- Service
- ConfigMap

Configuration of fields in Deployment such as number of replicas or resource limit and requirement will be supported in the future. 

#### Physical network topology [WIP]
By configuring the deployment strategies, Physical network topology of the benchmark microservices application can be altered. Features such as Service Mesh type, pod affinity, replicas, and load balancers can be configured via Kubernetes resources.

### Anomaly simulation [WIP]
the-bench extends chaos engineering tools to simulate more complex and realistic set of anomalies, faults, and failures in microservices.

## Project Structure
- `./cmd/service-unit/main.go` is the entry point for service unit binary.
- `./cmd/tbctl/main.go` is the entry point for cli program that can be used to validate the service-unit configuration and generate the Kubernetes manifests. 
- `./example/` holds the example configs.

## How it works
see [internals](./docs/internals.md).

