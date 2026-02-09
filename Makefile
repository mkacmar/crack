.PHONY: build build-release test test-unit test-e2e test-e2e-coverage clean lint fmt install-tools doc

BINARY = crack
ENTRYPOINT = ./cmd/crack
DIST_DIR = dist
COVERAGE_DIR = coverage

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS = -X github.com/mkacmar/crack/internal/version.Version=$(VERSION) \
          -X github.com/mkacmar/crack/internal/version.GitCommit=$(COMMIT) \
          -X github.com/mkacmar/crack/internal/version.BuildTime=$(BUILD_TIME)
RELEASE_LDFLAGS = -s -w $(LDFLAGS)

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(ENTRYPOINT)

build-release: lint test
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/$(BINARY)-linux-amd64 $(ENTRYPOINT)
	GOOS=linux GOARCH=arm64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/$(BINARY)-linux-arm64 $(ENTRYPOINT)
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/$(BINARY)-darwin-amd64 $(ENTRYPOINT)
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/$(BINARY)-darwin-arm64 $(ENTRYPOINT)
	GOOS=windows GOARCH=amd64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/$(BINARY)-windows-amd64.exe $(ENTRYPOINT)
	@echo ""
	@echo "SHA256 Checksums:"
	@cd $(DIST_DIR) && sha256sum *

test:
	go test -v ./...

test-unit:
	go test -v ./internal/...

test-e2e: build
	go test -v ./test/e2e/...

test-e2e-coverage:
	@mkdir -p $(COVERAGE_DIR)/raw
	go build -cover -ldflags "$(LDFLAGS)" -o $(BINARY) $(ENTRYPOINT)
	GOCOVERDIR=$(shell pwd)/$(COVERAGE_DIR)/raw go test -v ./test/e2e/...
	go tool covdata textfmt -i=$(COVERAGE_DIR)/raw -o=$(COVERAGE_DIR)/e2e.out
	go tool cover -html=$(COVERAGE_DIR)/e2e.out -o $(COVERAGE_DIR)/e2e.html
	xdg-open $(COVERAGE_DIR)/e2e.html

clean:
	rm -f $(BINARY)
	rm -rf $(DIST_DIR)
	rm -rf $(COVERAGE_DIR)

lint:
	go vet ./...
	staticcheck ./...

fmt:
	goimports -l -w .

install-tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

doc:
	go run ./internal/tools/doc
