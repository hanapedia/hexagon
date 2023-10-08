package config

import (
	"os"
	"strconv"
	"sync"

	model "github.com/hanapedia/the-bench/pkg/api/v1"
)

type EnvVars struct {
	DEP_ENV                  string
	TRACING                  bool
	HTTP_PORT                string
	GRPC_PORT                string
	KAFKA_PORT               string
	KAFKA_CLUSTER_NAME       string
	KAFKA_CLUSTER_NAMESPACE  string
	MONGO_USER               string
	MONGO_PASSWORD           string
	MONGO_PORT               string
	REDIS_PORT               string
	POSTGRE_PORT             string
	OTEL_COLLECTOR_NAME      string
	OTEL_COLLECTOR_NAMESPACE string
	OTEL_COLLECTOR_PORT      string
}

var envVars *EnvVars
var once sync.Once

func GetEnvs() *EnvVars {
	once.Do(func() {
		envVars = loadEnvVariables()
	})
	return envVars
}

func loadEnvVariables() *EnvVars {
	return &EnvVars{
		DEP_ENV:                  readEnv("DEP_ENV", "k8s"),
		TRACING:                  readBoolEnv("TRACING", false),
		HTTP_PORT:                readEnv("HTTP_PORT", strconv.Itoa(model.HTTP_PORT)),
		GRPC_PORT:                readEnv("GRPC_PORT", strconv.Itoa(model.GRPC_PORT)),
		KAFKA_PORT:               readEnv("KAFKA_PORT", strconv.Itoa(model.KAFKA_PORT)),
		KAFKA_CLUSTER_NAME:       readEnv("KAFKA_CLUSTER_NAME", model.KAFKA_CLUSTER_NAME),
		KAFKA_CLUSTER_NAMESPACE:  readEnv("KAFKA_CLUSTER_NAMESPACE", model.KAFKA_NAMESPACE),
		MONGO_USER:               readEnv("MONGO_USER", model.MONGO_USERNAME),
		MONGO_PASSWORD:           readEnv("MONGO_PASSWORD", model.MONGO_PASSWORD),
		MONGO_PORT:               readEnv("MONGO_PORT", strconv.Itoa(model.MONGO_PORT)),
		REDIS_PORT:               readEnv("REDIS_PORT", strconv.Itoa(model.REDIS_PORT)),
		POSTGRE_PORT:             readEnv("POSTGRE_PORT", strconv.Itoa(model.POSTGRES_PORT)),
		OTEL_COLLECTOR_NAME:      readEnv("OTEL_COLLECTOR_NAME", model.OTEL_COLLECTOR_NAME),
		OTEL_COLLECTOR_NAMESPACE: readEnv("OTEL_COLLECTOR_NAMESPACE", model.OTEL_COLLECTOR_NAMESPACE),
		OTEL_COLLECTOR_PORT:      readEnv("OTEL_COLLECTOR_PORT", strconv.Itoa(model.OTEL_COLLECTOR_PORT)),
	}
}

func readEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func readBoolEnv(key string, defaultValue bool) bool {
	boolValue := defaultValue
	if value, ok := os.LookupEnv(key); ok {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return boolValue
		}
		return parsed
	}
	return boolValue
}
