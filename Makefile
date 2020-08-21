NAME := ct
VERSION := $$(make -s show-version)
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS := "-s -w -X main.revision=$(CURRENT_REVISION)"
GOBIN ?= $(shell go env GOPATH)/bin

export GO111MODULE=on

.PHONY: all
all: clean build

.PHONY: build
build: bin-data
	CGO_ENABLED=1 go build -ldflags=$(BUILD_LDFLAGS) -o $(NAME) .

.PHONY: bin-data
bin-data:
	cd db/migrations && go-bindata -o migrations.go -pkg migrations -ignore migrations.go . && cd ../../

.PHONY: clean
clean:
	rm -rf $(NAME) goxz
	go clean

.PHONY: install
install:
	go install -ldflags=$(BUILD_LDFLAGS) .

.PHONY: tests
tests: clean build
	@echo "+ $@"
	go test ./...

.PHONY: integration-tests
integration-tests: clean install
	@echo "+ $@"
	bats -t tests/integration/*.bats

.PHONY: all-tests
all-tests: clean tests integration-tests
	@echo "+ $@"

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
