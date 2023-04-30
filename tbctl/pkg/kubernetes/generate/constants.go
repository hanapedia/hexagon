package generate

// temploral constants. will be replaced
const (
	NAMESPACE          = "default"
	SERVICE_UNIT_IMAGE = "hiroki11hanada/service-unit:v1.0"
	MONGO_IMAGE        = "hiroki11hanada/stateful-unit-mongo:v1.0"
	POSTGRE_IMAGE      = "hiroki11hanada/stateful-unit-postgre:v1.0"
	REPLICAS           = 1
	LIMIT_CPU          = "125m"
	LIMIT_MEM          = "128Mi"
	REQUEST_CPU        = "250m"
	REQUEST_MEM        = "64Mi"
	HTTP_PORT          = 8080
	GRPC_PORT          = 9090
	MONGO_PORT         = 27017
	POSTGRE_PORT       = 5432
	KAFKA_CLUSTER_NAME = "kafka-cluster"
	KAFKA_NAMESPACE    = "kafka-systems"
	KAFKA_PARTITIONS   = 3
	KAFKA_REPLICATIONS = 3
)
