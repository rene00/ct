NAME := ct
PKG := $(NAME)
VERSIONPKG := $(PKG)/pkg/version

# Set any default go build tags
BUILDTAGS :=

# Add to compile time flags
GITCOMMIT := $(shell git rev-parse --short HEAD)
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
BUILDTIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
ifneq ($(GITUNTRACKEDCHANGES),)
	GITCOMMIT := $(GITCOMMIT)-dirty
endif
GITBRANCH = $(shell git rev-parse --verify --abbrev-ref HEAD)
CTIMEVAR = -X '$(VERSIONPKG).commitSHA=$(GITCOMMIT)' \
		   -X '$(VERSIONPKG).branch=$(GITBRANCH)' \
		   -X '$(VERSIONPKG).date=$(BUILDTIME)'

GO_LDFLAGS=-ldflags "-w $(CTIMEVAR)"
GO_LDFLAGS_STATIC=-ldflags "-w $(CTIMEVAR) -extldflags -static"

all: clean

$(NAME): ## Builds a static executable
	@echo "+ $@"
	CGO_ENABLED=1 go build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(NAME)

.PHONY: build
build: $(NAME)

.PHONY: clean
clean:
	@echo "+ $@"
	$(RM) $(NAME)

.PHONY: tests
tests: build
	@echo "+ $@"
	go test ./... -cover

.PHONY: integration-tests
integration-tests: clean build
	@echo "+ $@"
	bats -t tests/integration/*.bats

.PHONY: all-tests
all-tests: clean tests integration-tests
	@echo "+ $@"

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
