package constants

const (
	// Name of the cluster cofig file
	CLUSTER_CONFIG_FILE_NAME = "hexagon-cluster.yaml"
	// Name of the config file mounted to the service unit pod
	SERVICE_UNIT_CONFIG_FILE_NAME = "service-unit.yaml"
)

type SeverAdapterVariant string
type BrokerVariant string
type RepositoryVariant string
type StressorValiant string

const (
	REST SeverAdapterVariant = "rest"
	GRPC SeverAdapterVariant = "grpc"

	KAFKA    BrokerVariant = "kafka"
	RABBITMQ BrokerVariant = "rabbitmq"
	Pulsar   BrokerVariant = "pulsar"

	MONGO   RepositoryVariant = "mongo"
	REDIS   RepositoryVariant = "redis"
	POSTGRE RepositoryVariant = "postgre"

	CPU    StressorValiant = "cpu"
	MEMORY StressorValiant = "memory"
	DISK   StressorValiant = "disk"
)

type Action string

const (
	READ          Action = "read"
	WRITE         Action = "write"
	GET           Action = "get"
	POST          Action = "post"
	SIMPLE_RPC    Action = "simpleRpc"
	CLIENT_STREAM Action = "clientStream"
	SERVER_STREAM Action = "serverStream"
	BI_STREAM     Action = "biStream"
)

type HttpMethod string

const (
	HTTP_POST HttpMethod = "POST"
	HTTP_GET  HttpMethod = "GET"
)

type PayloadSizeVariant string

const (
	SMALL  PayloadSizeVariant = "small"  // 1kb entries
	MEDIUM PayloadSizeVariant = "medium" // 4kb entries
	LARGE  PayloadSizeVariant = "large"  // 16kb entries
)

var PayloadSizeMap = map[PayloadSizeVariant]int64{
	SMALL:  1024,
	MEDIUM: 4096,
	LARGE:  16384,
}

const (
	DefaultPayloadSize = MEDIUM
)

const (
	DefaultPayloadCount = 3
)

const (
	NumInitialEntries = 100
)

const (
	DISK_STRESSOR_TMP_FILEPATH = "/tmp/shared_io_file"
)
