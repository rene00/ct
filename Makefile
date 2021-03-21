BIN := ct
VERSION := $$(make -s show-version)
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS := "-s -w -X main.revision=$(CURRENT_REVISION)"
GOBIN ?= $(shell go env GOPATH)/bin
export GO111MODULE=on

.PHONY: all
all: build

ct:
	CGO_ENABLED=1 go build -ldflags=$(BUILD_LDFLAGS) -o $(BIN) ./cmd/ct

.PHONY: build
build: clean bin-data ct

$(GOBIN)/go-bindata:
	cd && go get github.com/go-bindata/go-bindata/...

.PHONY: bin-data
bin-data: $(GOBIN)/go-bindata
	cd db/migrations && go-bindata -o bindata.go -pkg migrations -ignore bindata.go . && cd ../../

.PHONY: clean
clean:
	rm -rf $(BIN) goxz
	go clean

.PHONY: install
install: build
	mv -f ./ct $(GOBIN)/ct

.PHONY: test
test: clean build
	go test -cover ./...

$(GOBIN)/golint:
	cd && go get golang.org/x/lint/golint

$(GOBIN)/goxz:
	cd && go get github.com/Songmu/goxz/cmd/goxz

.PHONY: cross
cross: $(GOBIN)/goxz
	goxz -n $(BIN) -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) .

PHONY: show-version
show-version: $(GOBIN)/gobump
	@gobump show -r cmd/ct

$(GOBIN)/gobump:
	@cd && go get github.com/x-motemen/gobump/cmd/gobump

.PHONY: lint
lint: $(GOBIN)/golint
	go vet ct/...
	golint -set_exit_status internal/...

.PHONY: upload
upload: $(GOBIN)/ghr
	ghr "v$(VERSION)" goxz

$(GOBIN)/ghr:
	cd && go get github.com/tcnksm/ghr 

.PHONY: bump
bump: $(GOBIN)/gobump
ifneq ($(shell git status --porcelain),)
	$(error git workspace is dirty)
endif
ifneq ($(shell git rev-parse --abbrev-ref HEAD),master)
	$(error current branch is not master)
endif
	@gobump up -w cmd/ct
	git commit -am "bump up version to $(VERSION)"
	git tag "v$(VERSION)"
	git push origin master
	git push origin "refs/tags/v$(VERSION)"

