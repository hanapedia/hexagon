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
	return fmt.Sprintf("kafka:%s", port)
}

func GetMongoConnectionUri(adapterConfig model.StatefulAdapterConfig) string {
	port := GetEnvs().MONGO_PORT
	return fmt.Sprintf("mongodb://root:password@%s:%s/mongo?authSource=admin", adapterConfig.Name, port)
}
