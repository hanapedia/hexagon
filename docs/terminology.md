# Terminology
- Service Unit: refers to individual microservice in the benchmark application.
- Ingress Adapter: an interface on service unit for serving requests from other service units, or handling messages from brokers. For e.g.
    - REST API route
    - Kafka Consumer on a topic
- Egress Adapter: an interface on service unit for making requests to other service units, sending messages to brokers, or reading/writing to stateful services. e.g.
    - HTTP client calling a route
    - Kafka Producer for a topic
    - MongoDB client for a collection

Each service unit have one or more ingress adapters and each ingress adapter have ordered set of egress adapters. 
Each egress adapter will invoke ingress adapter on another service unit, send message to a broker, or perform transaction on a stateful service unit.
When a service unit receives a request/message to one of its ingress adapters, egress adapters associated to the ingress adapter will be called sequentially or in parallel.
