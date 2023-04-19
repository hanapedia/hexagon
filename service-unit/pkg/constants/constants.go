package constants

const (
	ServiceNameIndex = 0
	ProtocolIndex    = 1
	ActionIndex      = 2
	AdapterNameIndex = 3
)

type AdapterProtocol string

const (
	REST AdapterProtocol = "rest"
	// GRPC ServerAdapterProtocol = "grpc"
	KAFKA AdapterProtocol = "kafka"
)

const (
	PayloadSize = 16
)

const (
	RestServerAddr  = ":8080"
	KafkaBrokerAddr = "kafka:9092"
)
