# APP INFO
BUILD_DIR         := $(shell mktemp -d)
DOCKER_REGISTRY   ?= quay.io
IMAGE_PREFIX      ?= airshipit
IMAGE_NAME        ?= kubernetes-entrypoint
IMAGE_TAG         ?= latest
PROXY             ?= http://proxy.foo.com:8000
NO_PROXY          ?= localhost,127.0.0.1,.svc.cluster.local
USE_PROXY         ?= false
PUSH_IMAGE        ?= false
# use this variable for image labels added in internal build process
LABEL             ?= org.airshipit.build=community
COMMIT            ?= $(shell git rev-parse HEAD)
PYTHON            = python3
CHARTS            := $(filter-out deps, $(patsubst charts/%/.,%,$(wildcard charts/*/.)))
DISTRO             ?= ubuntu_noble
DISTRO_ALIAS	   ?= ubuntu_noble
IMAGE             := ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/${IMAGE_NAME}:${IMAGE_TAG}-${DISTRO}
IMAGE_ALIAS              := ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/${IMAGE_NAME}:${IMAGE_TAG}-${DISTRO_ALIAS}
UBUNTU_BASE_IMAGE ?=

# VERSION INFO
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

ifdef VERSION
	DOCKER_VERSION = $(VERSION)
endif

SHELL = /bin/bash

info:
	@echo "Version:           ${VERSION}"
	@echo "Git Tag:           ${GIT_TAG}"
	@echo "Git Commit:        ${GIT_COMMIT}"
	@echo "Git Tree State:    ${GIT_DIRTY}"
	@echo "Docker Version:    ${DOCKER_VERSION}"
	@echo "Registry:          ${DOCKER_REGISTRY}"
# Image URL to use all building/pushing image targets
IMG ?= quay.io/airshipit/kubernetes-entrypoint:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif


_BASE_IMAGE_ARG := $(if $(UBUNTU_BASE_IMAGE),--build-arg FROM="${UBUNTU_BASE_IMAGE}" ,)

MAKE_TARGET := build

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# linting
LINTER_CMD          := "github.com/golangci/golangci-lint/cmd/golangci-lint" run
LINTER_CONFIG       := .golangci.yaml


PKG                 := ./...
TESTS               := .

.PHONY: all
all: build unit-tests

.PHONY: build
build:
	@mkdir -p bin
	@CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o bin/kubernetes-entrypoint

.PHONY: lint
lint: build
	@go run ${LINTER_CMD} --config ${LINTER_CONFIG}

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: docker-image
docker-image:
ifeq ($(USE_PROXY), true)
	@docker build --network host -t $(IMAGE) --label $(LABEL) \
		--label "org.opencontainers.image.revision=$(COMMIT)" \
		--label "org.opencontainers.image.created=$(shell date --rfc-3339=seconds --utc)" \
		--label "org.opencontainers.image.title=$(IMAGE_NAME)" \
		-f images/Dockerfile.$(DISTRO) \
		$(_BASE_IMAGE_ARG) \
		--build-arg MAKE_TARGET=$(MAKE_TARGET) \
		--build-arg http_proxy=$(PROXY) \
		--build-arg https_proxy=$(PROXY) \
		--build-arg HTTP_PROXY=$(PROXY) \
		--build-arg HTTPS_PROXY=$(PROXY) \
		--build-arg no_proxy=$(NO_PROXY) \
		--build-arg NO_PROXY=$(NO_PROXY) .
else
	@docker build --network host -t $(IMAGE) --label $(LABEL) \
		--label "org.opencontainers.image.revision=$(COMMIT)" \
		--label "org.opencontainers.image.created=$(shell date --rfc-3339=seconds --utc)" \
		--label "org.opencontainers.image.title=$(IMAGE_NAME)" \
		-f images/Dockerfile.$(DISTRO) \
		--build-arg MAKE_TARGET=$(MAKE_TARGET) \
		$(_BASE_IMAGE_ARG) .
endif
ifneq ($(DISTRO), $(DISTRO_ALIAS))
	docker tag $(IMAGE) $(IMAGE_ALIAS)
ifeq ($(DOCKER_REGISTRY), localhost:5000)
	docker push $(IMAGE_ALIAS)
endif
endif
ifeq ($(DOCKER_REGISTRY), localhost:5000)
	docker push $(IMAGE)
endif
ifeq ($(PUSH_IMAGE), true)
	docker push $(IMAGE)
endif

check-docker:
	@if [ -z $$(which docker) ]; then \
		echo "Missing \`docker\` client which is required for development"; \
		exit 2; \
	fi

images: check-docker docker-image

.PHONY: docker-image-unit-tests
docker-image-unit-tests: MAKE_TARGET = unit-tests
docker-image-unit-tests: docker-image

.PHONY: docker-image-lint
docker-image-lint: MAKE_TARGET = lint
docker-image-lint: docker-image


.PHONY: get-modules
get-modules:
	@go mod download

.PHONY: clean
clean:
	@rm -rf bin

.PHONY: unit-test
unit-tests: build
	@go test -v ./...
