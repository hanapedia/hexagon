package config

import (
	"os"
	"strconv"
	"sync"
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
		HTTP_PORT:                readEnv("HTTP_PORT", "8080"),
		GRPC_PORT:                readEnv("GRPC_PORT", "9090"),
		KAFKA_PORT:               readEnv("KAFKA_PORT", "9092"),
		KAFKA_CLUSTER_NAME:       readEnv("KAFKA_CLUSTER_NAME", "my-cluster"),
		KAFKA_CLUSTER_NAMESPACE:  readEnv("KAFKA_CLUSTER_NAMESPACE", "kafka"),
		MONGO_USER:               readEnv("MONGO_USER", "root"),
		MONGO_PASSWORD:           readEnv("MONGO_PASSWORD", "password"),
		MONGO_PORT:               readEnv("MONGO_PORT", "27017"),
		REDIS_PORT:               readEnv("REDIS_PORT", "6379"),
		POSTGRE_PORT:             readEnv("POSTGRE_PORT", "5432"),
		OTEL_COLLECTOR_NAME:      readEnv("OTEL_COLLECTOR_NAME", "otelcollector-collector"),
		OTEL_COLLECTOR_NAMESPACE: readEnv("OTEL_COLLECTOR_NAMESPACE", "observability"),
		OTEL_COLLECTOR_PORT:      readEnv("OTEL_COLLECTOR_PORT", "4317"),
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
