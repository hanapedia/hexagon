package factory

const (
	NAMESPACE            = "default"
	SERVICE_UNIT_IMAGE   = "hiroki11hanada/service-unit:v1.5"
	MONGO_IMAGE          = "hiroki11hanada/stateful-unit-mongo:v1.0"
	LOAD_GENERATOR_IMAGE = "hiroki11hanada/tb-load-generator:v1.5"
	// POSTGRE_IMAGE      = "hiroki11hanada/stateful-unit-postgre:v1.0"
	REPLICAS     = 1
	REQUEST_CPU  = "125m"
	REQUEST_MEM  = "64Mi"
	LIMIT_CPU    = "250m"
	LIMIT_MEM    = "128Mi"
	LIMIT_MEM_LG = "1Gi"
	HTTP_PORT    = 8080
	// GRPC_PORT          = 9090
	MONGO_PORT = 27017
	// POSTGRE_PORT       = 5432
	KAFKA_CLUSTER_NAME = "my-cluster" // default name from strimzi kafka operator getting started guide
	KAFKA_NAMESPACE    = "kafka"      // default name from strimzi kafka operator getting started guide
	KAFKA_PARTITIONS   = 1
	KAFKA_REPLICATIONS = 1
)
