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
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/tbctl cmd/tbctl/main.go
