# Variables
DOCKER_HUB_USERNAME := hiroki11hanada
REPOSITORY_NAME := tb-load-generator
TAG ?= v0.0.2
ENGINE ?= podman
IMAGE_NAME := $(DOCKER_HUB_USERNAME)/$(REPOSITORY_NAME):$(TAG)

# Conditionally set the BUILD_COMMAND and PUSH_COMMAND
ifeq ($(ENGINE),docker)
	BUILD_COMMAND=docker build -t $(IMAGE_NAME) .
	PUSH_COMMAND=docker push $(IMAGE_NAME)
else ifeq ($(ENGINE),podman)
	BUILD_COMMAND=podman build -t $(IMAGE_NAME) .
	PUSH_COMMAND=podman push $(IMAGE_NAME)
else
$(error ENGINE must be either 'docker' or 'podman')
endif

.PHONY: build
build:
	$(BUILD_COMMAND)

.PHONY: push
push: build
	$(PUSH_COMMAND)

.PHONY: all
all: build push
