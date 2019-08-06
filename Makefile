GOOS          ?= $(shell go env GOOS)
GOARCH        ?= $(shell go env GOARCH)

.PHONY: build
build:
	@echo "Building kubernetes-entrypoint for $(GOOS)/$(GOARCH)"
	mkdir -p bin/$(GOARCH)
	go build -o bin/$(GOARCH)/kubernetes_entrypoint

.PHONY: get-modules
get-modules:
	@go mod download

.PHONY: linux-arm64
linux-arm64:
linux-arm64: GOOS = "linux"
linux-arm64: GOARCH = "arm64"
linux-arm64: build

.PHONY: linux-amd64
linux-amd64: GOOS = "linux"
linux-amd64: GOARCH = "amd64"
linux-amd64: build

.PHONY: clean
clean:
	rm -rf bin

.PHONY: test
test:
	go test ./...

all: linux-amd64 linux-arm64
