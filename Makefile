# Binary names
BINARY=nes-outage-status-checker
HEALTHCHECK_BINARY=healthcheck

# Build flags
LDFLAGS=-s -w

.PHONY: build build-healthcheck build-all clean run deps fmt vet

## build: Build the main binary
build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) .

## build-healthcheck: Build the health check server
build-healthcheck:
	go build -ldflags="$(LDFLAGS)" -o $(HEALTHCHECK_BINARY) ./cmd/healthcheck

## build-all: Build all binaries
build-all: build build-healthcheck

## clean: Remove built binaries
clean:
	rm -f $(BINARY) $(HEALTHCHECK_BINARY)

## deps: Download dependencies
deps:
	go mod download

## fmt: Format code
fmt:
	go fmt ./...

## vet: Run go vet
vet:
	go vet ./...

## help: Show this help
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
