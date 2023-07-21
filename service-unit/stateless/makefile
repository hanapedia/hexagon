# Variables
DOCKER_HUB_USERNAME := hiroki11hanada
REPOSITORY_NAME := service-unit
TAG ?= v0.0.2
ENGINE ?= podman
IMAGE_NAME := $(DOCKER_HUB_USERNAME)/$(REPOSITORY_NAME):$(TAG)

# Conditionally set the BUILD_COMMAND and PUSH_COMMAND
ifeq ($(ENGINE),docker)
	BUILD_COMMAND=docker buildx build --platform linux/amd64 -t $(IMAGE_NAME) . --load
	BUILD_COMMAND_ARM64=docker buildx build --platform linux/arm64 -t $(IMAGE_NAME) . --load
	BUILD_COMMAND_ALL=docker buildx build --platform linux/amd64,linux/arm64 -t $(IMAGE_NAME) . --push
	PUSH_COMMAND=docker push $(IMAGE_NAME)
else ifeq ($(ENGINE),podman)
	BUILD_COMMAND=podman build --platform linux/amd64 -t $(IMAGE_NAME) .
	BUILD_COMMAND_ARM64=podman build --platform linux/arm64 -t $(IMAGE_NAME) .
	BUILD_COMMAND_ALL=podman build --platform linux/amd64,linux/arm64 -t $(IMAGE_NAME) .
	PUSH_COMMAND=podman push $(IMAGE_NAME)
else
$(error ENGINE must be either 'docker' or 'podman')
endif

.PHONY: dev
dev:
	go mod edit -dropreplace=github.com/hanapedia/the-bench/the-bench-operator
	go mod tidy
	$(BUILD_COMMAND_ARM64)
	go mod edit -replace=github.com/hanapedia/the-bench/the-bench-operator=../../the-bench-operator
	$(PUSH_COMMAND)

.PHONY: release
release:
	go mod edit -dropreplace=github.com/hanapedia/the-bench/the-bench-operator
	go mod tidy
	$(BUILD_COMMAND)
	go mod edit -replace=github.com/hanapedia/the-bench/the-bench-operator=../../the-bench-operator
	$(PUSH_COMMAND)

.PHONY: push
push:
	$(PUSH_COMMAND)

.PHONY: all
all:
	go mod edit -dropreplace=github.com/hanapedia/the-bench/the-bench-operator
	go mod tidy
	$(BUILD_COMMAND_ALL)
	go mod edit -replace=github.com/hanapedia/the-bench/the-bench-operator=../../the-bench-operator
	$(PUSH_COMMAND)
