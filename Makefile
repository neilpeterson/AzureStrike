.PHONY: build run test test-cover clean fmt vet lint help

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
	$(GO) test -cover ./...

## test-verbose: Run tests with verbose output
test-verbose:
	$(GO) test -v ./...

## fmt: Format code
fmt:
	$(GO) fmt ./...

## vet: Run go vet
vet:
	$(GO) vet ./...

## clean: Remove build artifacts
clean:
	rm -f $(BINARY_NAME)
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
