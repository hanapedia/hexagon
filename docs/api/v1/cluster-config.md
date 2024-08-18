## Cluster Config
These configurations applies globally to all service units. 
The configuration must be provided in a file named `hexagon-cluster.yaml` at the directory specified via `hexctl generate -f`.

Most of these configuration values can also be provided via environmental variable through [deployment configuration](./deployment-config.md), and be overridden by individual service unit config.

Default values are set to all of the parameters and everything should work just fine without specifying any of these parameters.

The definitions for Go struct is in [`pkg/api/v1/clusterconfig.go`](../../../pkg/api/v1/clusterconfig.go)

| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| namespace     | Kubernetes Namespace to deploy the generated application. | "hexagon" | false |
| logLevel      | Log Level for each service unit. | "info" | false |
| dockerHubUsername | Container registry for the image used by the service units. | "hexagonbenchamrk" | false |
| healthPort   | Port used to start health server. | 6060 | false |
| metricsPort   | Port used to start metris server. | 7070 | false |
| httpPort | Port used to start HTTP REST server. | 8080 | false |
| grpcPort | Port used to start gRPC server. | 9090 | false |
| tracing.enabled | Whether tracing is enabled. | false | false |
| kafka | Configuration values for Kafka cluster | default values in [#Kafka](#kafka) | false |
| mongo | Configuration values for Mongo | default values in [#Mongo](#mongo) | false |
| redis | Configuration values for Redis | default values in [#Redis](#redis) | false ||
| otel  | Configuration values for Otel Collector | default values in [#Otel Collector](#otel-collector) | false |

### Kafka
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| namespace     | Kubernetes Namespace in which Kafka cluster is running | "kafka" | false |
| name          | Name of the Kafka cluster deployed by Strimzi. | "my-cluster" | false |
| port          | Port used by the Brokers. | 9092 | false |
| replications  | Replications setting. | 1 | false |
| partitions    | Partitions setting. | 1 | false |

### Mongo
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| imageName     | Image name of Mongo container | "stateful-unit-mongo" | false |
| port          | Port used by Mongo. | 27017 | false |
| username      | Username for admin user. | "root" | false |
| password      | Password for admin user. | "password" | false |

### Redis
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| imageName     | Image name of Redis container | "stateful-unit-redis" | false |
| port          | Port used by Redis. | 6379 | false |

### Otel Collector
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| namespace     | Kubernetes Namespace in which otel collector is running | "monitoring" | false |
| name          | Name of the otel collector. | "otel-collector-collector" | false |
| port          | Port used by otel collector. | 4317 | false |
