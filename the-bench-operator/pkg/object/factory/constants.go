package factory

const (
	NAMESPACE = "the-bench"

	SERVICE_UNIT_IMAGE_NAME   = "hiroki11hanada/service-unit"
	MONGO_IMAGE_NAME          = "hiroki11hanada/stateful-unit-mongo"
	LOAD_GENERATOR_IMAGE_NAME = "hiroki11hanada/tb-load-generator"

	REPLICAS     = 1
	REQUEST_CPU  = "125m"
	REQUEST_MEM  = "64Mi"
	LIMIT_CPU    = "250m"
	LIMIT_MEM    = "128Mi"
	LIMIT_MEM_LG = "1Gi"

	HTTP_PORT  = 8080
	MONGO_PORT = 27017

	KAFKA_CLUSTER_NAME = "my-cluster" // default name from strimzi kafka operator getting started guide
	KAFKA_NAMESPACE    = "kafka"      // default name from strimzi kafka operator getting started guide
	KAFKA_PARTITIONS   = 1
	KAFKA_REPLICATIONS = 1

	CHAOSMESH_NAMESPACE      = "chaos-mesh"
	CHAOSMESH_DURATION       = "5m"
	CHAOSMESH_LATENCY        = "15ms"
	CHAOSMESH_LATENCY_JITTER = "5ms"
)
