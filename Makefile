BIN := ct
VERSION := $$(make -s show-version)
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS := "-s -w -X main.revision=$(CURRENT_REVISION)"
GOBIN ?= $(shell go env GOPATH)/bin

export GO111MODULE=on

.PHONY: all
all: clean build

.PHONY: build
build: bin-data
	CGO_ENABLED=1 go build -ldflags=$(BUILD_LDFLAGS) -o $(BIN) .

$(GOBIN)/go-bindata:
	cd && go get github.com/go-bindata/go-bindata/...

.PHONY: bin-data
bin-data: $(GOBIN)/go-bindata
	cd db/migrations && go-bindata -o migrations.go -pkg migrations -ignore migrations.go . && cd ../../

.PHONY: clean
clean:
	rm -rf $(BIN) goxz
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

$(GOBIN)/golint:
	cd && go get golang.org/x/lint/golint

$(GOBIN)/goxz:
	cd && go get github.com/Songmu/goxz/cmd/goxz

.PHONY: cross
cross: $(GOBIN)/goxz
	goxz -n $(BIN) -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) .

PHONY: show-version
show-version: $(GOBIN)/gobump
	@gobump show -r .

$(GOBIN)/gobump:
	@cd && go get github.com/x-motemen/gobump/cmd/gobump

.PHONY: lint
lint: $(GOBIN)/golint
	go vet .
	golint -set_exit_status . cmd config internal/...


.PHONY: upload
upload: $(GOBIN)/ghr
	ghr "v$(VERSION)" goxz

$(GOBIN)/ghr:
	cd && go get github.com/tcnksm/ghr 
