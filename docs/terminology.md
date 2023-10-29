# Terminology
The name Hexagon comes from Hexagonal Architecture, which is a design pattern commonly used for developing a microservice application. The terms used comes from the concepts defined in the architecture. 

- Service Unit: refers to individual microservice in the benchmark application.
- Primary Adapter: an interface on service unit for serving requests from other service units, or handling messages from brokers. For e.g.
    - REST API route
    - Kafka Consumer on a topic
- Secondary Adapter: an interface on service unit for making requests to other service units, sending messages to brokers, or reading/writing to stateful services. e.g.
    - HTTP client calling a route
    - Kafka Producer for a topic
    - MongoDB client for a collection

## How they work together
- Each service unit have one or more primary adapters and each primary adapter have ordered set of secondary adapters. 
- Each secondary adapter will invoke primary adapter on another service unit, send message to a broker, or perform transaction on a stateful service unit.
- When a service unit receives a request/message to one of its primary adapters, secondary adapters associated to the primary adapter will be called sequentially or in parallel.
