package constants

type SeverAdapterVariant string
type BrokerVariant string
type RepositoryVariant string

const (
	REST SeverAdapterVariant = "rest"
	GRPC SeverAdapterVariant = "grpc"

	KAFKA    BrokerVariant = "kafka"
	RABBITMQ BrokerVariant = "rabbitmq"
	Pulsar   BrokerVariant = "pulsar"

	MONGO   RepositoryVariant = "mongo"
	POSTGRE RepositoryVariant = "postgre"
)

type Action string

const (
	READ  Action = "read"
	WRITE Action = "write"
)

type HttpMethod string

const (
	POST HttpMethod = "POST"
	GET  HttpMethod = "GET"
)

type PayloadSizeVariant string
type PayloadSize int

const (
	SMALL      PayloadSizeVariant = "small"  // 1kb entries
	MEDIUM     PayloadSizeVariant = "medium" // 4kb entries
	LARGE      PayloadSizeVariant = "large"  // 16kb entries
	SMALLSIZE  PayloadSize    = 1        // 1kb entries
	MEDIUMSIZE PayloadSize    = 4        // 1kb entries
	LARGESIZE  PayloadSize    = 16       // 1kb entries
)

const (
	DefaultPayloadSize = MEDIUMSIZE
)

const (
	NumInitialEntries = 100
)
