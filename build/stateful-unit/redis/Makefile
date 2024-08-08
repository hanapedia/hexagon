# Variables
DOCKER_HUB_USERNAME := hexagonbenchmark
REPOSITORY_NAME := stateful-unit-redis
TAG ?= v0.0.3
IMAGE_NAME := $(DOCKER_HUB_USERNAME)/$(REPOSITORY_NAME):$(TAG)
BUILDER_NAME := datagen
BUILDER_TAG ?= $(TAG)
BUILDER_IMAGE := $(DOCKER_HUB_USERNAME)/$(BUILDER_NAME):$(BUILDER_TAG)

	BUILD_COMMAND=docker buildx build --build-arg BUILDER_IMAGE=$(BUILDER_IMAGE) --platform linux/amd64 -t $(IMAGE_NAME) . --load
	BUILD_COMMAND_ARM64=docker buildx build --build-arg BUILDER_IMAGE=$(BUILDER_IMAGE) --platform linux/arm64 -t $(IMAGE_NAME) . --load
	PUSH_COMMAND=docker push $(IMAGE_NAME)

.PHONY: dev prod push
dev:
	$(BUILD_COMMAND_ARM64)
	$(PUSH_COMMAND)

prod:
	$(BUILD_COMMAND)
	$(PUSH_COMMAND)

push:
	$(PUSH_COMMAND)
