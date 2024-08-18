package config

import (
	"os"
	"strconv"
	"sync"

	"github.com/hanapedia/hexagon/pkg/api/defaults"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/utils"
	corev1 "k8s.io/api/core/v1"
)

type EnvVars struct {
	LOG_LEVEL                string
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
	OTEL_COLLECTOR_NAME      string
	OTEL_COLLECTOR_NAMESPACE string
	OTEL_COLLECTOR_PORT      string
	METRICS_PORT             string
	HEALTH_PORT              string
}

var envVars *EnvVars
var once sync.Once

func FromClusterConfig(config model.ClusterConfig) *EnvVars {
	return &EnvVars{
		LOG_LEVEL:                config.LogLevel,
		TRACING:                  config.Tracing.Enabled,
		HTTP_PORT:                strconv.Itoa(int(config.HTTPPort)),
		GRPC_PORT:                strconv.Itoa(int(config.GRPCPort)),
		KAFKA_PORT:               strconv.Itoa(int(config.Kafka.Port)),
		KAFKA_CLUSTER_NAME:       config.Kafka.ClusterName,
		KAFKA_CLUSTER_NAMESPACE:  config.Kafka.Namespace,
		MONGO_USER:               config.Mongo.Username,
		MONGO_PASSWORD:           config.Mongo.Password,
		MONGO_PORT:               strconv.Itoa(int(config.Mongo.Port)),
		REDIS_PORT:               strconv.Itoa(int(config.Redis.Port)),
		OTEL_COLLECTOR_NAME:      config.Otel.Name,
		OTEL_COLLECTOR_NAMESPACE: config.Otel.Namespace,
		OTEL_COLLECTOR_PORT:      strconv.Itoa(int(config.Otel.Port)),
		METRICS_PORT:             strconv.Itoa(int(config.MetricsPort)),
		HEALTH_PORT:              strconv.Itoa(int(config.HealthPort)),
	}
}

func (e EnvVars) AsK8sEnvVars() []corev1.EnvVar {
	k8sEnvs := make([]corev1.EnvVar, 0)
	k8sEnvs = append(k8sEnvs, corev1.EnvVar{Name: "TRACING", Value: utils.Btos(e.TRACING)})
	if e.LOG_LEVEL != "" {
		k8sEnvs = append(k8sEnvs, corev1.EnvVar{Name: "LOG_LEVEL", Value: e.LOG_LEVEL})
	}
	if e.HTTP_PORT != "0" {
		k8sEnvs = append(k8sEnvs, corev1.EnvVar{Name: "HTTP_PORT", Value: e.HTTP_PORT})
	}
	if e.GRPC_PORT != "0" {
		k8sEnvs = append(k8sEnvs, corev1.EnvVar{Name: "GRPC_PORT", Value: e.GRPC_PORT})
	}
	if e.KAFKA_PORT != "0" {
		k8sEnvs = append(k8sEnvs, corev1.EnvVar{Name: "KAFKA_PORT", Value: e.KAFKA_PORT})
	}
	if e.KAFKA_CLUSTER_NAME != "" {
		k8sEnvs = append(k8sEnvs, corev1.EnvVar{Name: "KAFKA_CLUSTER_NAME", Value: e.KAFKA_CLUSTER_NAME})
	}
	if e.KAFKA_CLUSTER_NAMESPACE != "" {
		k8sEnvs = append(k8sEnvs, corev1.EnvVar{Name: "KAFKA_CLUSTER_NAMESPACE", Value: e.KAFKA_CLUSTER_NAMESPACE})
	}
	if e.MONGO_USER != "" {
		k8sEnvs = append(k8sEnvs, corev1.EnvVar{Name: "MONGO_USER", Value: e.MONGO_USER})
	}
	if e.MONGO_PASSWORD != "" {
		k8sEnvs = append(k8sEnvs, corev1.EnvVar{Name: "MONGO_PASSWORD", Value: e.MONGO_PASSWORD})
	}
	if e.MONGO_PORT != "0" {
		k8sEnvs = append(k8sEnvs, corev1.EnvVar{Name: "MONGO_PORT", Value: e.MONGO_PORT})
	}
	if e.REDIS_PORT != "0" {
		k8sEnvs = append(k8sEnvs, corev1.EnvVar{Name: "REDIS_PORT", Value: e.REDIS_PORT})
	}
	if e.OTEL_COLLECTOR_NAME != "" {
		k8sEnvs = append(k8sEnvs, corev1.EnvVar{Name: "OTEL_COLLECTOR_NAME", Value: e.OTEL_COLLECTOR_NAME})
	}
	if e.OTEL_COLLECTOR_NAMESPACE != "" {
		k8sEnvs = append(k8sEnvs, corev1.EnvVar{Name: "OTEL_COLLECTOR_NAMESPACE", Value: e.OTEL_COLLECTOR_NAMESPACE})
	}
	if e.OTEL_COLLECTOR_PORT != "0" {
		k8sEnvs = append(k8sEnvs, corev1.EnvVar{Name: "OTEL_COLLECTOR_PORT", Value: e.OTEL_COLLECTOR_PORT})
	}

	return k8sEnvs
}

func GetEnvs() *EnvVars {
	once.Do(func() {
		envVars = loadEnvVariables()
	})
	return envVars
}

func loadEnvVariables() *EnvVars {
	return &EnvVars{
		LOG_LEVEL:                readEnv("LOG_LEVEL", defaults.LOG_LEVEL),
		TRACING:                  readBoolEnv("TRACING", defaults.TRACING),
		HTTP_PORT:                readEnv("HTTP_PORT", strconv.Itoa(defaults.HTTP_PORT)),
		GRPC_PORT:                readEnv("GRPC_PORT", strconv.Itoa(defaults.GRPC_PORT)),
		KAFKA_PORT:               readEnv("KAFKA_PORT", strconv.Itoa(defaults.KAFKA_PORT)),
		KAFKA_CLUSTER_NAME:       readEnv("KAFKA_CLUSTER_NAME", defaults.KAFKA_CLUSTER_NAME),
		KAFKA_CLUSTER_NAMESPACE:  readEnv("KAFKA_CLUSTER_NAMESPACE", defaults.KAFKA_NAMESPACE),
		MONGO_USER:               readEnv("MONGO_USER", defaults.MONGO_USERNAME),
		MONGO_PASSWORD:           readEnv("MONGO_PASSWORD", defaults.MONGO_PASSWORD),
		MONGO_PORT:               readEnv("MONGO_PORT", strconv.Itoa(defaults.MONGO_PORT)),
		REDIS_PORT:               readEnv("REDIS_PORT", strconv.Itoa(defaults.REDIS_PORT)),
		OTEL_COLLECTOR_NAME:      readEnv("OTEL_COLLECTOR_NAME", defaults.OTEL_COLLECTOR_NAME),
		OTEL_COLLECTOR_NAMESPACE: readEnv("OTEL_COLLECTOR_NAMESPACE", defaults.OTEL_COLLECTOR_NAMESPACE),
		OTEL_COLLECTOR_PORT:      readEnv("OTEL_COLLECTOR_PORT", strconv.Itoa(defaults.OTEL_COLLECTOR_PORT)),
		METRICS_PORT:             readEnv("METRICS_PORT", strconv.Itoa(defaults.METRICS_PORT)),
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
