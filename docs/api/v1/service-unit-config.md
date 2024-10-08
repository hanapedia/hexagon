## Service Unit Config
The definitions for Go struct is in:
- service unit: [`pkg/api/v1/serviceunitconfig.go`](../../../pkg/api/v1/serviceunitconfig.go)
- primary adapter: [`pkg/api/v1/primaryadapter.go`](../../../pkg/api/v1/primaryadapter.go)
- secondary adapter: [`pkg/api/v1/secondaryadapter.go`](../../../pkg/api/v1/secondaryadapter.go)

| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| name          | Name of the service.                          | ""          | true        |
| version       | Version of the image used.                    | ""          | true        |
| deployment    | [Deployment configs](./deployment-config.md). | {}          | false       |
| adapters      | List of [Primary Adapter Configs](#primary-adapter).| []    | true        |

### Primary Adapter
Primary adapter can be type of server, repository, or cosumer and only one of the configuration should be given.
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| server        | Configuration for [server](#server).          | {}          | false       |
| repository    | Configuration for [repository](#repository).  | {}          | false       |
| consumer      | Configuration for [consumer](#consumer).      | {}          | false       |
| tasks         | List of [tasks](#task) (secondary adapters) attached to primary adapter. | [] | false |

### Server
Server configuration is defined in the unit of individual routes in REST APi.
Each definition of server primary adapter corresponds to a route attached to REST API.

*For gRPC server, while gRPC does not have the concept of routes, the unit of configuration is the same, 
meaning that each definition will represent independent execution of secondary adapters.*
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| variant       | Variant of the server. "rest" or "grpc"       | ""          | true        |
| action        | Action for the route. "post" or "get" for rest and "simpleRpc", "clientStream", "serverStream", or "biStream" for grpc | "" | true |
| route         | Unique identifier for the route.               | ""          | true        |
| payload       | Size of the payload that this route returns. Can be "small", "medium", or "large" | "medium" | false |
| weight        | Weight for the route to be called when load generator is enabled for the service. If load generator is not enabled, the field will be ignored | 0 | false |
| payloadCount  | Number of payloads to return for grpc route with serverStream action. | 3 | false |

### Repository
Repository configuration indicates that the service is a stateful service, which uses different image from regular service units.

*When this configuration is given, any other primary adapter configuration in the adapters list will be ignored.*
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| variant       | Variant of the repository. "mongo" or "redis" | ""          | true        |

### Consumer
Consumer configuration is defined in the unit of topic that the service unit cosumes. 

*When this configuration is given, the topic manifest for strimzi kafka will be generated*
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| variant       | Variant of the consumer. Only "kafka" is supported at the moment. | "" | true |
| topic         | Name of the topic that the consumer subscribes to. | ""     | true        |

### Task
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| adapter | Configuration for the [Secondary Adatper](#secondary-adapter). | {} | true |
| concurrent | Whether to execute the task concurrently. | false      | false        |
| resiliency | [Resiliency patterns](./resiliency-config.md) to apply for the secondary adapter. | {}      | false        |

### Secondary Adapter
Secondary adapter can be type of invocation, repository, producer, or stressor and only one of the configuration should be given.
Each secondary adapter should match existing primary adapters on other services.
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| invocation    | Configuration for [invocation](#invocation). This will invoke server primary adapters on other services. | {} | false |
| repository    | Configuration for [repository](#repository-client). This will read from or write to stateful services. | {} | false |
| producer      | Configuration for [producer](#producer). This will produce message to specified topic. | {} | false |
| stressor      | Configuration for [stressor](#stressor). This will create internal stress within the service. | {} | false |

#### Invocation
Invocation configuration specify which server primary adapter on other service is called.
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| service       | Name of the service to invoke.                | ""          | true        |
| variant       | Variant of the server. "rest" or "grpc"       | ""          | true        |
| action        | Action for the route. "post" or "get" for rest and "simpleRpc", "clientStream", "serverStream", or "biStream" for grpc | "" | true |
| route         | Unique identifier for the route.               | ""          | true        |
| payload       | Size of the payload that this call sends. Can be "small", "medium", or "large" | "medium" | false |

#### Repository Client
Repository client configuration specify which repository primary adapter on other service is called.
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| service       | Name of the service to invoke.                | ""          | true        |
| variant       | Variant of the repository. "mongo" or "redis" | ""          | true        |
| action        | Action for the route. "read" or "write"       | ""          | true        |
| payload       | Size of the payload that this call sends. Can be "small", "medium", or "large" | "medium" | false |

#### Producer
Producer configuration specify which topic to send the message to. 
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| variant       | Variant of the broker. only "kafka" is supported | ""       | true        |
| topic         | Name of the topic to publish message.         | ""          | true        |
| payload       | Size of the payload that this call sends. Can be "small", "medium", or "large" | "medium" | false |

#### Stressor
Stressor configuration specify how the service is stressed internally.
This is treated the same way as other secondary adapter, except it does not send any external requests.
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| name          | Unique name assigned for the stressor.        | ""          | true        |
| variant       | Variant of the stressor. only "cpu" is supported. | ""      | true        |
| duration      | Duration with units.                          | ""          | true        |
| threadCount   | Number of threads to spawn.                   | ""          | false       |
| payload       | Size of the payload that this call sends. Can be "small", "medium", or "large" | "medium" | false |
