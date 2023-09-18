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

type RepositoryEntryVariant string
type RepositoryEntrySize int

const (
	SMALL      RepositoryEntryVariant = "small"  // 1kb entries
	MEDIUM     RepositoryEntryVariant = "medium" // 4kb entries
	LARGE      RepositoryEntryVariant = "large"  // 16kb entries
	SMALLSIZE  RepositoryEntrySize    = 1        // 1kb entries
	MEDIUMSIZE RepositoryEntrySize    = 4        // 1kb entries
	LARGESIZE  RepositoryEntrySize    = 16       // 1kb entries
)

const (
	PayloadSize = LARGESIZE
)

const (
	NumInitialEntries = 100
)
