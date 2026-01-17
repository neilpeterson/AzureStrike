.PHONY: build run test test-cover test-race test-verbose clean fmt vet lint ci help

BINARY_NAME=azurestrike
BUILD_DIR=.
GO=go

## build: Build the binary
build:
	$(GO) build -o $(BINARY_NAME) ./cmd/azurestrike

## run: Run in development mode
run:
	$(GO) run ./cmd/azurestrike

## run-scenario: Run a specific scenario (usage: make run-scenario SCENARIO=storage-breach)
run-scenario:
	./$(BINARY_NAME) --scenario $(SCENARIO)

## list: List available scenarios
list:
	./$(BINARY_NAME) --list

## test: Run all tests
test:
	$(GO) test ./...

## test-cover: Run tests with coverage report
test-cover:
	$(GO) test -v -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -func=coverage.out

## test-race: Run tests with race detector
test-race:
	$(GO) test -race ./...

## test-verbose: Run tests with verbose output
test-verbose:
	$(GO) test -v ./...

## fmt: Format code
fmt:
	$(GO) fmt ./...

## vet: Run go vet
vet:
	$(GO) vet ./...

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## ci: Run all CI checks locally
ci: fmt vet lint test-race test-cover

## clean: Remove build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	$(GO) clean

## deps: Download dependencies
deps:
	$(GO) mod download

## tidy: Tidy go modules
tidy:
	$(GO) mod tidy

## help: Show this help message
help:
	@echo "AzureStrike - Azure Security Training Game"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'
