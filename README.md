# Hexagon
A highly configurable microservices benchmark application generator designed specifically for researchers.
Hexagon can generate benchmark microservices applications of virtually unlimited scale and interaction patterns without requiring application rebuilds.

![Hexagon Diagram](./docs/assets/hexagon_main_diagram.png)

## Quick Start
Follow this guide to test out Hexagon.
This guide will walk you through basic steps for intallation and manifest generation.

### Requirements

| Dependency    | Version    | Description                            |
| ------------- | ---------- | -------------------------------------- |
| Go            | 1.22       | Should also work with 1.21, 1.20       |
| Kind          | 1.27+      | Should work with any standard Kubernetes cluster |

### 1. Build the CLI
Currently precompiled binary is not available but will be ready in the future.
For now, you need to compile the CLI after cloning this repository.
```sh
go build -o ./hexctl cmd/hexctl/main.go
```

### 2. Start Kind Cluster
You can skip this step if you already have a working cluster.
```sh
kind create cluster --name hexagon-cluster
```

### 3. Generate Manifests
Generate manifest with the CLI built ealier.
This will generate manifests for example application that emulates [Online Boutique](https://github.com/GoogleCloudPlatform/microservices-demo), using the example configuration available in example/config/onlineboutique.
For more details on configuration values, please refer to the [documentaion](./docs/configuration.md).
```sh
./hexctl generate -f example/config/onlineboutique/ -o example/manifest/
```
`-f` specifies the directory for the hexagon configuration files
`-o` specifies the directory for the output Kubernetes manifets

### 4. Apply the manifest
Create namespace and apply the manifests, nothing special.
```sh
# create namespace
kubectl create namespace hexagon
# hexagon generates manifests for `hexagon` namespace by default
kubectl apply -f example/manifest/
```

### 5. (Optional) Run load generator
There is a ready to use K6 loadgenerator deployment for the emulated application available in `example/loadgenerator`.
Apply it to see the generated application handling traffic.
```sh
kubectl apply -n hexagon -f example/loadgenerator/manifests.yaml
```

Currently, Hexagon cannot generate UIs so watch logs to see the application running.
For example, Watch the logs for `frontend`.
```sh
kubectl logs -n hexagon --selector app=frontend -f
```

### Clean up
Simply delete the Kind cluster if you created one.
```sh
kind delete cluster --name hexagon-cluster
```

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

