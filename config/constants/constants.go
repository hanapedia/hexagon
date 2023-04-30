package constants

type StatelessAdapterVariant string
type BrokerAdapterVariant string
type StatefulAdapterVariant string
const (
	REST StatelessAdapterVariant = "rest"
	GRPC StatelessAdapterVariant = "grpc"

	KAFKA    BrokerAdapterVariant = "kafka"
	RABBITMQ BrokerAdapterVariant = "rabbitmq"
	Pulsar   BrokerAdapterVariant = "pulsar"

	MONGO   StatefulAdapterVariant = "mongo"
	POSTGRE StatefulAdapterVariant = "postgre"
)

type Action string
const (
	READ  Action = "read"
	WRITE Action = "write"
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
