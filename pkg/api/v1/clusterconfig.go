package v1

import (
	"time"

	"github.com/hanapedia/hexagon/pkg/api/defaults"
)

type ClusterConfig struct {
	ConfigTemplate
	Namespace         string                     `json:"namespace,omitempty"`
	LogLevel          string                     `json:"logLevel,omitempty"`
	ServiceMonitor    bool                       `json:"serviceMonitor,omitempty"`
	DockerHubUsername string                     `json:"dockerHubUsername,omitempty"`
	MetricsPort       int32                      `json:"metricsPort,omitempty"`
	HealthPort        int32                      `json:"healthPort,omitempty"`
	HTTPPort          int32                      `json:"httpPort,omitempty"`
	GRPCPort          int32                      `json:"grpcPort,omitempty"`
	Tracing           TracingClusterConfig       `json:"tracing,omitempty"`
	Kafka             KafkaClusterConfig         `json:"kafka,omitempty"`
	Mongo             MongoClusterConfig         `json:"mongo,omitempty"`
	Redis             RedisClusterConfig         `json:"redis,omitempty"`
	Otel              OtelCollectorClusterConfig `json:"otel,omitempty"`
	Resiliency        ResiliencySpec             `json:"resiliency,omitempty"`
	BaseTimeout       string                     `json:"baseTimeout,omitempty"`
}

type TracingClusterConfig struct {
	Enabled bool `json:"enabled,omitempty"`
}

type KafkaClusterConfig struct {
	Port         int32  `json:"port,omitempty"`
	ClusterName  string `json:"clusterName,omitempty"`
	Namespace    string `json:"namespace,omitempty"`
	Partitions   int32  `json:"partitions,omitempty"`
	Replications int32  `json:"replications,omitempty"`
}

type MongoClusterConfig struct {
	ImageName string `json:"imageName,omitempty"`
	Port      int32  `json:"port,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
}

type RedisClusterConfig struct {
	ImageName string `json:"imageName,omitempty"`
	Port      int32  `json:"port,omitempty"`
}

type OtelCollectorClusterConfig struct {
	Port      int32  `json:"port,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

func NewClusterConfig() ClusterConfig {
	return ClusterConfig{
		Namespace:         defaults.NAMESPACE,
		LogLevel:          defaults.LOG_LEVEL,
		ServiceMonitor:    defaults.SERVICE_MONITOR,
		DockerHubUsername: defaults.DOCKER_USERNAME,
		MetricsPort:       defaults.METRICS_PORT,
		HealthPort:        defaults.HEALTH_PORT,
		HTTPPort:          defaults.HTTP_PORT,
		GRPCPort:          defaults.GRPC_PORT,
		Tracing:           TracingClusterConfig{Enabled: defaults.TRACING},
		Kafka: KafkaClusterConfig{
			Port:         defaults.KAFKA_PORT,
			ClusterName:  defaults.KAFKA_CLUSTER_NAME,
			Namespace:    defaults.KAFKA_NAMESPACE,
			Partitions:   defaults.KAFKA_PARTITIONS,
			Replications: defaults.KAFKA_REPLICATIONS,
		},
		Mongo: MongoClusterConfig{
			Port:      defaults.MONGO_PORT,
			Username:  defaults.MONGO_USERNAME,
			Password:  defaults.MONGO_PASSWORD,
			ImageName: defaults.MONGO_IMAGE_NAME,
		},
		Redis: RedisClusterConfig{
			Port:      defaults.REDIS_PORT,
			ImageName: defaults.REDIS_IMAGE_NAME,
		},
		Otel: OtelCollectorClusterConfig{
			Port:      defaults.OTEL_COLLECTOR_PORT,
			Name:      defaults.OTEL_COLLECTOR_NAME,
			Namespace: defaults.OTEL_COLLECTOR_NAMESPACE,
		},
	}
}

// Parses BaseTimeout with time.ParseDuration
func (cc *ClusterConfig) GetBaseTimeout() (time.Duration, error) {
	duration, err := time.ParseDuration(cc.BaseTimeout)
	if err != nil {
		return duration, err
	}
	return duration, nil
}
