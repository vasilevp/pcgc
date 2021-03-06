export GO111MODULE := on
export PATH := ./bin:$(PATH)

TEST?=$$(go list ./...)
GOFMT_FILES?=$$(find . -name '*.go')
GOLANGCI_VERSION=v1.22.2

default: build

.PHONY: setup
setup:
	@echo "==> Installing dependencies..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLANGCI_VERSION)
	go get -u github.com/google/addlicense

.PHONY: test
test:
	@echo "==> Running tests..."
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

.PHONY: lint
lint:
	@echo "==> Linting all packages..."
	golangci-lint run ./... -E gofmt -E golint -E misspell

.PHONY: fmt
fmt:
	gofmt -s -w $(GOFMT_FILES); goimports -w $(GOFMT_FILES)

.PHONY: addlicense
addlicense:
	addlicense -c "MongoDB Inc" $(GOFMT_FILES)

.PHONY: test-compile
test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./package-name"; \
		exit 1; \
	fi
	@echo "==> Compiling test binary..."
	go test -c $(TEST) $(TESTARGS)

# Build targets
.PHONY: clean
clean:
	@echo "==> Cleaning build artifacts..."
	go clean ./...

gitsha := $(shell git log -n1 --pretty='%h')
version=$(shell git describe --exact-match --tags "$(gitsha)" 2>/dev/null)
ifeq ($(version),)
	version := $(gitsha)
endif
ldflags=-ldflags='-X github.com/mongodb-labs/pcgc/pkg/httpclient.version=$(version)'
.PHONY: build
build:
	go build $(ldflags) ./...

# GIT hooks
.PHONY: link-git-hooks
link-git-hooks:
	@echo "==> Installing all git hooks..."
	find .git/hooks -type l -exec rm {} \;
	find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;
