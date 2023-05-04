package config

import (
	"fmt"

	"github.com/hanapedia/the-bench/config/model"
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

func GetMongoConnectionUri(adapterConfig model.StatefulAdapterConfig) string {
	port := GetEnvs().MONGO_PORT
	return fmt.Sprintf("mongodb://root:password@%s:%s/mongo?authSource=admin", adapterConfig.Name, port)
}
