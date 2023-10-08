GRPC_GENERATE_DIR := ./internal/service-unit/infrastructure/adapters/generated/grpc/
GO_MODULE := $(shell go list -m)
COMMA := ,

# Local development with tilt
.PHONY: devstart
devstart:
	ctlptl apply -f ./dev/cluster.yaml

	# install strimzi kafka operator and sample kafka cluster
	kubectl create namespace kafka
	kubectl create -n kafka -f 'https://strimzi.io/install/latest?namespace=kafka'
	kubectl -n kafka wait --for=condition=available --timeout=180s --all deployments
	kubectl apply -n kafka -f https://strimzi.io/examples/latest/kafka/kafka-persistent-single.yaml
	kubectl wait -n kafka kafka/my-cluster --for=condition=Ready --timeout=300s

	# create the-bench namespace
	kubectl create namespace the-bench

.PHONY: devstop
devstop:
	ctlptl delete -f ./dev/cluster.yaml

.PHONY: devmanifests
devmanifests:
	./bin/tbctl generate -f ./dev/config/ -o ./dev/manifest/

.PHONY: devbuild
devbuild:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/service-unit cmd/service-unit/main.go

.PHONY: devbuildcli
devbuildcli:
	CGO_ENABLED=0 go build -o bin/tbctl cmd/tbctl/main.go

.PHONY: devbuilddatagen
devbuilddatagen:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/datagen cmd/datagen/main.go

.PHONY: genproto
genproto:
	protoc --go_out=$(GRPC_GENERATE_DIR) \
		--go_opt=paths=source_relative$(COMMA)Mproto/service-unit.proto=$(GO_MODULE)/$(GRPC_GENERATE_DIR) \
		--go-grpc_out=$(GRPC_GENERATE_DIR) \
		--go-grpc_opt=paths=source_relative$(COMMA)Mproto/service-unit.proto=$(GO_MODULE)/$(GRPC_GENERATE_DIR) \
		./proto/service-unit.proto
