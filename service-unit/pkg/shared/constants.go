package shared

const (
 ServiceNameIndex = 0
 ProtocolIndex = 1
 ActionIndex = 2
 AdapterNameIndex = 3

)

type ServerAdapterProtocol string

const (
  REST ServerAdapterProtocol = "rest"
  // GRPC ServerAdapterProtocol = "grpc"
  // KAFKA ServerAdapterProtocol = "kafka"
)
