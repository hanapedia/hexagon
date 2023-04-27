# Variables
DOCKER_HUB_USERNAME := hiroki11hanada
REPOSITORY_NAME := service-unit
TAG ?= v1.0
IMAGE_NAME := $(DOCKER_HUB_USERNAME)/$(REPOSITORY_NAME):$(TAG)

# Targets
.PHONY: build
build:
	go mod edit -dropreplace=github.com/hanapedia/the-bench/config
	go mod tidy
	docker build -t $(IMAGE_NAME) .
	go mod edit -replace=github.com/hanapedia/the-bench/config=../../config

.PHONY: push
push: build
	docker push $(IMAGE_NAME)

.PHONY: all
all: build push

