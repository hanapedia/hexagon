## Quick Start
Follow this guide to test out Hexagon.
This guide will walk you through basic tasks for intallation and manifest generation.
At the end of this guide, you will have a Hexagon generated version of Online Boutique microservices application.

### Requirements

| Dependency    | Version    | Description                            |
| ------------- | ---------- | -------------------------------------- |
| [Go](https://go.dev/doc/install) | 1.22       | Should also work with 1.21, 1.20       | 
| [kind](https://kind.sigs.k8s.io/docs/user/quick-start/) | 1.27+      | Should work with any standard Kubernetes cluster |

### 1. Build the CLI
Currently precompiled binary is not available but will be ready in the future.
For now, you need to compile the CLI after cloning this repository. 
```sh
go build -o ./hexctl cmd/hexctl/main.go
```

### 2. Start kind Cluster
You can skip this task if you already have a working cluster.
```sh
kind create cluster --name hexagon-cluster
```

### 3. Generate Manifests
Generate manifest with the CLI built ealier.
This will generate manifests for example application that emulates [Online Boutique](https://github.com/GoogleCloudPlatform/microservices-demo), using the example configuration available in example/onlineboutique/config/.
For more details on configuration values, please refer to the [documentaion](./docs/api/v1/).
```sh
./hexctl generate -f example/onlineboutique/config/ -o example/onlineboutique/manifest/
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
Simply delete the kind cluster if you created one.
```sh
kind delete cluster --name hexagon-cluster
```

### What's next?
Check out other examples in the [examples directory](../example/) to see different types of adapter definitions.

Also, see [configuration specs](./api/v1/) for more details.

Also, see [`hexctl` docs](./hexctl.md) for available commands for hexctl.

