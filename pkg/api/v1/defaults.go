package v1

const (
	NAMESPACE = "the-bench"

	SERVICE_UNIT_IMAGE_NAME   = "hiroki11hanada/service-unit"
	MONGO_IMAGE_NAME          = "hiroki11hanada/stateful-unit-mongo"
	REDIS_IMAGE_NAME          = "hiroki11hanada/stateful-unit-redis"
	LOAD_GENERATOR_IMAGE_NAME = "hiroki11hanada/tb-load-generator"

	REPLICAS     = 1
	REQUEST_CPU  = "125m"
	REQUEST_MEM  = "64Mi"
	LIMIT_CPU    = "250m"
	LIMIT_MEM    = "128Mi"
	LIMIT_MEM_LG = "1Gi"

	HTTP_PORT           = 8080
	GRPC_PORT           = 9090
	KAFKA_PORT          = 9092
	MONGO_PORT          = 27017
	REDIS_PORT          = 6379
	POSTGRES_PORT       = 5432
	OTEL_COLLECTOR_PORT = 4317

	KAFKA_CLUSTER_NAME = "my-cluster" // default name from strimzi kafka operator getting started guide
	KAFKA_NAMESPACE    = "kafka"      // default name from strimzi kafka operator getting started guide
	KAFKA_PARTITIONS   = 1
	KAFKA_REPLICATIONS = 1

	MONGO_USERNAME = "root"
	MONGO_PASSWORD = "password"

	OTEL_COLLECTOR_NAME      = "otelcollector-collector"
	OTEL_COLLECTOR_NAMESPACE = "observability"

	CHAOSMESH_NAMESPACE      = "chaos-mesh"
	CHAOSMESH_DURATION       = "5m"
	CHAOSMESH_LATENCY        = "15ms"
	CHAOSMESH_LATENCY_JITTER = "5ms"
)
