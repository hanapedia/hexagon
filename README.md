# the-bench
A highly configurable microservices benchmark suite.

## Description
the-bench is a microservices benchmark suite with configurablity in mind. It aims to solve the limited benchmarking possibilites for researchers studying microservices. Using the-bench researchers can easily construct a benchmark microservices application with highly configured compositions. the-bench also includes toolchain to simulate dynamic microservices application update.

## Features
### Application generation
the-bench allows configuration of following features for generated application
#### Logical network topology
By decralatively configuring each service's API and handlers, benchmark microservices application of any logical network topology can be generated. Each service can expose set of APIs such as REST, gRPC, or Kafka consumer, and attach multiple handlers for each APIs. Also each handler can access multiple APIs exposed by other services. All this interactions between services can be defined declaratively in a simple yaml configuration, without writing a single line of code.
* only REST API supported at the moment.
#### Physical network topology [WIP]
By configuring the deployment strategies, Physical network topology of the benchmark microservices application can be altered. Features such as Service Mesh type, pod affinity, replicas, and load balancers can be configured via Kubernetes resources.
#### Service Internal Workload
By configuring internal workload of each service, it can simulate different types of tasks that real world microservices application have. The internal workload will produce artificial stress for computing resources such as cpu, memory and disk.
#### Stateful Service [WIP]
the-bench allows the configuration of type and number of database used through out the application.

### Dynamic update simulation [WIP]
the-bench allows the redefinition of application component without taking down the entire application nor building additional images. The the-bench controller accepts and validates change in configuration at runtime, identify the services that require configuration changes, then applies the new configuration to subjected services. This process can also be scheduled for better control in experiments.

### Anomaly simulation [WIP]
the-bench extends chaos engineering tools to simulate more complex and realistic set of anomalies, faults, and failures in microservices.
