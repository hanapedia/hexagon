# Variables
DOCKER_HUB_USERNAME := hiroki11hanada
REPOSITORY_NAME := service-unit
TAG ?= v0.0.3
IMAGE_NAME := $(DOCKER_HUB_USERNAME)/$(REPOSITORY_NAME):$(TAG)

	BUILD_COMMAND=docker buildx build --platform linux/amd64 -t $(IMAGE_NAME) . --load
	BUILD_COMMAND_ARM64=docker buildx build --platform linux/arm64 -t $(IMAGE_NAME) . --load
	BUILD_COMMAND_ALL=docker buildx build --platform linux/amd64,linux/arm64 -t $(IMAGE_NAME) . --push
	PUSH_COMMAND=docker push $(IMAGE_NAME)

.PHONY: dev prod push all
dev:
	go mod edit -dropreplace=github.com/hanapedia/the-bench/the-bench-operator
	go mod tidy
	$(BUILD_COMMAND_ARM64)
	go mod edit -replace=github.com/hanapedia/the-bench/the-bench-operator=../../the-bench-operator
	$(PUSH_COMMAND)

prod:
	go mod edit -dropreplace=github.com/hanapedia/the-bench/the-bench-operator
	go mod tidy
	$(BUILD_COMMAND)
	go mod edit -replace=github.com/hanapedia/the-bench/the-bench-operator=../../the-bench-operator
	$(PUSH_COMMAND)

push:
	$(PUSH_COMMAND)

all:
	go mod edit -dropreplace=github.com/hanapedia/the-bench/the-bench-operator
	go mod tidy
	$(BUILD_COMMAND_ALL)
	go mod edit -replace=github.com/hanapedia/the-bench/the-bench-operator=../../the-bench-operator
	$(PUSH_COMMAND)
