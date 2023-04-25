# Variables
DOCKER_HUB_USERNAME := hiroki11hanada
REPOSITORY_NAME := stateful-unit-mongo
TAG ?= v1.0
IMAGE_NAME := $(DOCKER_HUB_USERNAME)/$(REPOSITORY_NAME):$(TAG)

# Targets

.PHONY: build
build:
	docker build -t $(IMAGE_NAME) .

.PHONY: push
push: build
	docker push $(IMAGE_NAME)

.PHONY: all
all: build push

