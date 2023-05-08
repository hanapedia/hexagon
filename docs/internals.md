# Internals
## Stateless Service Unit
Stateless servcie unit is implemented in a single docker image written in Go.
1. the container process created from the image reads the configuration file given in yaml once it is started.
2. then it instantiates the ingress and egress adapters accordingly.

## Stateful Service Unit
Custom docker image based on respective database images with initial data inserted.

## Broker
Kafka is deployed using [strimzi operator](https://strimzi.io/).
Topics are created when service unit with a broker adapter with the topic exist.
