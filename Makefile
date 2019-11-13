SHELL := /bin/bash

# linting
LINTER_CMD          := "github.com/golangci/golangci-lint/cmd/golangci-lint" run
LINTER_CONFIG       := .golangci.yaml

# docker image options
DOCKER_REGISTRY       ?= quay.io
DOCKER_IMAGE_NAME     ?= kubernetes-entrypoint
DOCKER_IMAGE_PREFIX   ?= airshipit
DOCKER_IMAGE_TAG      ?= dev
DOCKER_IMAGE          ?= $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_PREFIX)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
DOCKER_MAKE_TARGET    := build
DOCKER_TARGET_STAGE   ?= release
DOCKER_BASE_IMAGE     ?= docker.io/golang:1.12.6-stretch
DOCKER_RELEASE_IMAGE  ?= scratch

PKG                 := ./...
TESTS               := .

.PHONY: build
build:
	@mkdir -p bin
	@CGO_ENABLED=0 go build -o bin/kubernetes-entrypoint

.PHONY: lint
lint:
	@go run ${LINTER_CMD} --config ${LINTER_CONFIG}

.PHONY: docker-image
docker-image:
	@docker build . --tag $(DOCKER_IMAGE) --target $(DOCKER_TARGET_STAGE) \
		--build-arg MAKE_TARGET=$(DOCKER_MAKE_TARGET) \
		--build-arg GO_IMAGE=$(DOCKER_BASE_IMAGE) \
		--build-arg RELEASE_IMAGE=$(DOCKER_RELEASE_IMAGE)

.PHONY: docker-image-unit-tests
docker-image-unit-tests: DOCKER_MAKE_TARGET = unit-tests
docker-image-unit-tests: DOCKER_TARGET_STAGE = builder
docker-image-unit-tests: docker-image

.PHONY: docker-image-lint
docker-image-lint: DOCKER_MAKE_TARGET = lint
docker-image-lint: DOCKER_TARGET_STAGE = builder
docker-image-lint: docker-image

.PHONY: get-modules
get-modules:
	@go mod download

.PHONY: clean
clean:
	@rm -rf bin

.PHONY: unit-test
unit-tests:
	@go test -v ./...
