package config

import (
	"fmt"

	model "github.com/hanapedia/the-bench/pkg/api/v1"
)

func GetRestServerAddr() string {
	port := GetEnvs().HTTP_PORT
	return fmt.Sprintf(":%s", port)
}

func GetGrpcServerAddr() string {
	port := GetEnvs().GRPC_PORT
	return fmt.Sprintf(":%s", port)
}

func GetKafkaBrokerAddr() string {
	port := GetEnvs().KAFKA_PORT
	clusterName := GetEnvs().KAFKA_CLUSTER_NAME
	clusterNamespace := GetEnvs().KAFKA_CLUSTER_NAMESPACE
	depEnv := GetEnvs().DEP_ENV
	if depEnv == "docker" {
		return fmt.Sprintf("kafka:%s", port)
	}
	return fmt.Sprintf("%s-kafka-bootstrap.%s.svc.cluster.local:%s", clusterName, clusterNamespace, port)
}

func GetMongoConnectionUri(adapterConfig *model.RepositoryClientConfig) string {
	user := GetEnvs().MONGO_USER
	password := GetEnvs().MONGO_PASSWORD
	// port := GetEnvs().MONGO_PORT
	port := "27017" // port is hardcoded until environmental variable issue is resolved
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/mongo?authSource=admin", user, password, adapterConfig.Name, port)
}

func GetRedisConnectionAddr(adapterConfig *model.RepositoryClientConfig) string {
	port := GetEnvs().REDIS_PORT
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/mongo?authSource=admin", user, password, adapterConfig.Name, port)
}

func GetGrpcDialAddr(adapterConfig *model.InvocationConfig) string {
	port := GetEnvs().GRPC_PORT
	return fmt.Sprintf("%s:%s", adapterConfig.Service, port)
}

func GetOtelCollectorUrl() string {
	depEnv := GetEnvs().DEP_ENV
	if depEnv == "docker" {
		return fmt.Sprintf("otelcollector:%s", GetEnvs().OTEL_COLLECTOR_PORT)
	}
	return fmt.Sprintf(
		"%s.%s.svc.cluster.local:%s",
		GetEnvs().OTEL_COLLECTOR_NAME,
		GetEnvs().OTEL_COLLECTOR_NAMESPACE,
		GetEnvs().OTEL_COLLECTOR_PORT,
	)
}
