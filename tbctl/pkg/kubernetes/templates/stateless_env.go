package templates

const StatelessEnvs =
`HTTP_PORT={{ .HTTP_PORT }}
GRPC_PORT={{ .GRPC_PORT }}
KAFKA_PORT={{ .KAFKA_PORT }}
MONGO_PORT={{ .MONGO_PORT }}
POSTGRE_PORT={{ .POSTGRE_PORT }}`
