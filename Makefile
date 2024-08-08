GRPC_GENERATE_DIR := ./internal/service-unit/infrastructure/adapters/generated/grpc/
GO_MODULE := $(shell go list -m)
COMMA := ,
TEST_PATH ?= ./...
GO_VERSION := 1.22
INTEGRATION_TEST_DOCKERFILE_PATH := ./test/integration/Dockerfile
DOCKER_USER ?= hexagonbenchmark

.PHONY: ctestbuild
ctestbuild:
	docker build --build-arg WORKDIR=$(CURDIR) -t test_container -f $(INTEGRATION_TEST_DOCKERFILE_PATH) $(CURDIR)

# Containerized testing for local testing.
# This makes things easier by allowing port allocation for MacOS.
.PHONY: ctest
ctest: ctestbuild
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v $(CURDIR):$(CURDIR) -e TC_HOST=host.docker.internal test_container go test -v $(TEST_PATH)

# Alias for running tests.
.PHONY: test
test:
	go test -v $(TEST_PATH)

# Local development with tilt
#
.PHONY: devmini
devmini:
	ctlptl apply -f ./dev/cluster.yaml

	# create namespaces
	kubectl apply -f ./dev/namespaces.yaml

	# create curl pod
	kubectl apply -n hexagon -f ./dev/curl.yaml

#
.PHONY: devstart
devstart:
	ctlptl apply -f ./dev/cluster.yaml

	# create namespaces
	kubectl apply -f ./dev/namespaces.yaml

	# install strimzi kafka operator and sample kafka cluster
	kubectl apply -f https://raw.githubusercontent.com/hanapedia/lab-cluster-apps/main/kafka/operator/overlays/dev/manifests.yaml -n kafka
	kubectl -n kafka wait --for=condition=available --timeout=180s --all deployments

	# restart once
	kubectl rollout restart deployment -n kafka strimzi-cluster-operator
	kubectl -n kafka wait --for=condition=available --timeout=180s --all deployments

	sleep 10

	kubectl apply -f https://raw.githubusercontent.com/hanapedia/lab-cluster-apps/main/kafka/kafka/overlays/dev/manifests.yaml -n kafka
	kubectl wait -n kafka kafka/my-cluster --for=condition=Ready --timeout=300s

	# install tempo
	kubectl apply -f https://raw.githubusercontent.com/hanapedia/lab-cluster-apps/main/tempo/dev/manifests.yaml
	# install otel collector
	kubectl apply -f https://raw.githubusercontent.com/hanapedia/lab-cluster-apps/main/otel/collector/overlays/dev/manifests.yaml
	kubectl -n monitoring wait --for=condition=available --timeout=180s --all deployments

	# create curl pod
	kubectl apply -n hexagon -f ./dev/curl.yaml

.PHONY: devstop
devstop:
	ctlptl delete -f ./dev/cluster.yaml

.PHONY: devmanifests
devmanifests:
	bash ./dev/generate.sh "$(DOCKER_USER)"

.PHONY: devbuild
devbuild:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/service-unit cmd/service-unit/main.go

.PHONY: devbuildcli
devbuildcli:
	CGO_ENABLED=0 go build -o bin/hexctl cmd/hexctl/main.go

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
