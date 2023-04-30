package config

import (
	"os"
	"sync"
)

type EnvVars struct {
	HTTP_PORT    string
	GRPC_PORT    string
	KAFKA_PORT   string
	MONGO_PORT   string
	POSTGRE_PORT string
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
	httpPort := "8080"
	grpcPort := "9090"
	kafkaPort := "9092"
	mongoPort := "27017"
	postgrePort := "5432"

	if envHttpPort, ok := os.LookupEnv("HTTP_PORT"); ok {
		httpPort = envHttpPort
	}

	if envGrpcPort, ok := os.LookupEnv("GRPC_PORT"); ok {
		grpcPort = envGrpcPort
	}

	if envKafkaPort, ok := os.LookupEnv("KAFKA_PORT"); ok {
		kafkaPort = envKafkaPort
	}

	if envmongoPort, ok := os.LookupEnv("MONGO_PORT"); ok {
		mongoPort = envmongoPort
	}

	if envPostgrePort, ok := os.LookupEnv("POSTGRE_PORT"); ok {
		postgrePort = envPostgrePort
	}

	return &EnvVars{
		HTTP_PORT:    httpPort,
		GRPC_PORT:    grpcPort,
		KAFKA_PORT:   kafkaPort,
		MONGO_PORT:   mongoPort,
		POSTGRE_PORT: postgrePort,
	}
}
