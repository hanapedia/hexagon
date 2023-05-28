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
	depEnv := "k8s"
	tracing := false
	httpPort := "8080"
	grpcPort := "9090"
	kafkaPort := "9092"
	kafkaClusterName := "my-cluster"
	kafkaClusterNamespace := "kafka"
	mongoUser := "root"
	mongoPassword := "password"
	mongoPort := "27017"
	postgrePort := "5432"
	otelCollectorName := "otelcollector"
	otelCollectorNamespace := "observability"
	otelCollectorPort := "4317"

	if envDepEnv, ok := os.LookupEnv("DEP_ENV"); ok {
		depEnv = envDepEnv
	}

	if envTracing, ok := os.LookupEnv("TRACING"); ok {
		parsed, err := strconv.ParseBool(envTracing)
		if err != nil {
			tracing = false
		}
		tracing = parsed
	}

	if envHttpPort, ok := os.LookupEnv("HTTP_PORT"); ok {
		httpPort = envHttpPort
	}

	if envGrpcPort, ok := os.LookupEnv("GRPC_PORT"); ok {
		grpcPort = envGrpcPort
	}

	if envKafkaPort, ok := os.LookupEnv("KAFKA_PORT"); ok {
		kafkaPort = envKafkaPort
	}

	if envKafkaClusterName, ok := os.LookupEnv("KAFKA_PORT"); ok {
		kafkaClusterName = envKafkaClusterName
	}

	if envKafkaClusterNamespace, ok := os.LookupEnv("KAFKA_PORT"); ok {
		kafkaClusterNamespace = envKafkaClusterNamespace
	}

	if envmongoUser, ok := os.LookupEnv("MONGO_USER"); ok {
		mongoUser = envmongoUser
	}

	if envmongoPassword, ok := os.LookupEnv("MONGO_PASSWORD"); ok {
		mongoPassword = envmongoPassword
	}

	if envmongoPort, ok := os.LookupEnv("MONGO_PORT"); ok {
		mongoPort = envmongoPort
	}

	if envPostgrePort, ok := os.LookupEnv("POSTGRE_PORT"); ok {
		postgrePort = envPostgrePort
	}

	if envOtelCollectorName, ok := os.LookupEnv("OTEL_COLLECTOR_NAME"); ok {
		otelCollectorName = envOtelCollectorName
	}

	if envOtelCollectorNamespace, ok := os.LookupEnv("OTEL_COLLECTOR_NAMESPACE"); ok {
		otelCollectorNamespace = envOtelCollectorNamespace
	}

	if envOtelCollectorPort, ok := os.LookupEnv("OTEL_COLLECTOR_PORT"); ok {
		otelCollectorPort = envOtelCollectorPort
	}

	return &EnvVars{
		DEP_ENV:                  depEnv,
		TRACING:                  tracing,
		HTTP_PORT:                httpPort,
		GRPC_PORT:                grpcPort,
		KAFKA_PORT:               kafkaPort,
		KAFKA_CLUSTER_NAME:       kafkaClusterName,
		KAFKA_CLUSTER_NAMESPACE:  kafkaClusterNamespace,
		MONGO_USER:               mongoUser,
		MONGO_PASSWORD:           mongoPassword,
		MONGO_PORT:               mongoPort,
		POSTGRE_PORT:             postgrePort,
		OTEL_COLLECTOR_NAME:      otelCollectorName,
		OTEL_COLLECTOR_NAMESPACE: otelCollectorNamespace,
		OTEL_COLLECTOR_PORT:      otelCollectorPort,
	}
}
