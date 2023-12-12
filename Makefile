GRPC_GENERATE_DIR := ./internal/service-unit/infrastructure/adapters/generated/grpc/
GO_MODULE := $(shell go list -m)
COMMA := ,

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
	rm -f ./dev/manifest/all/generated/* && ./bin/hexctl generate -f ./dev/config/all -o ./dev/manifest/all/generated
	rm -f ./dev/manifest/rest/generated/* && ./bin/hexctl generate -f ./dev/config/rest -o ./dev/manifest/rest/generated/
	rm -f ./dev/manifest/kafka/generated/* && ./bin/hexctl generate -f ./dev/config/kafka -o ./dev/manifest/kafka/generated
	rm -f ./dev/manifest/redis/generated/* && ./bin/hexctl generate -f ./dev/config/redis -o ./dev/manifest/redis/generated
	rm -f ./dev/manifest/mongo/generated/* && ./bin/hexctl generate -f ./dev/config/mongo -o ./dev/manifest/mongo/generated
	rm -f ./dev/manifest/grpc/generated/* && ./bin/hexctl generate -f ./dev/config/grpc -o ./dev/manifest/grpc/generated/
	rm -f ./dev/manifest/onlineboutique/generated/* && ./bin/hexctl generate -f ./dev/config/onlineboutique -o ./dev/manifest/onlineboutique/generated/
	rm -f ./dev/manifest/onlineboutique-async/generated/* && ./bin/hexctl generate -f ./dev/config/onlineboutique-async -o ./dev/manifest/onlineboutique-async/generated/

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
