package config

import (
	"fmt"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
)

func GetRestServerAddr() string {
	port := GetEnvs().HTTP_PORT
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
	port := GetEnvs().MONGO_PORT
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/mongo?authSource=admin", user, password, adapterConfig.Name, port)
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
