BINARY  := pong
ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
VERSION := $(or $(IMAGE_TAG),$(shell git describe --tags --first-parent --match "v*" 2> /dev/null || echo v0.0.0))
LOCAL_BIN_DIR := $(ROOT_DIR)/bin

.DEFAULT_GOAL := help

.PHONY: $(BINARY) test coverage help html html-gen html-wsl

build: ## Build the binary
	go build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o $(BINARY) cmd/$(BINARY)/main.go

test: ## Run unit tests
	@go test -race ./...

coverage: ## Run unit tests with coverage
	@go test -v -race -cover -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out

html: ## Show html coverage result
	@go tool cover -html=./coverage.out

html-gen: ## Export html coverage result
	@go tool cover -html=./coverage.out -o ./coverage.html

html-wsl: html-gen ## Open html coverage result in wsl
	@explorer.exe `wslpath -w ./coverage.html` || true

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
