.PHONY: build build-release test test-unit test-e2e test-e2e-coverage clean lint fmt install-tools doc

BINARY = crack
ENTRYPOINT = ./cmd/crack
DIST_DIR = dist
COVERAGE_DIR = coverage

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
TIME_FORMAT = %Y-%m-%dT%H:%M:%SZ
BUILD_TIME ?= $(shell git log -1 --format=%cd --date=format:'$(TIME_FORMAT)' 2>/dev/null || date -u +'$(TIME_FORMAT)')
LDFLAGS = -X go.kacmar.sk/crack/internal/version.Version=$(VERSION) \
          -X go.kacmar.sk/crack/internal/version.GitCommit=$(COMMIT) \
          -X go.kacmar.sk/crack/internal/version.BuildTime=$(BUILD_TIME)
RELEASE_LDFLAGS = -s -w $(LDFLAGS)
GOFLAGS = -buildmode=pie

build:
	go build $(GOFLAGS) -tags debug -ldflags "$(LDFLAGS)" -o $(BINARY) $(ENTRYPOINT)

build-release: lint test
	@mkdir -p $(DIST_DIR)/$(VERSION)
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/$(VERSION)/$(BINARY)_$(VERSION)_linux_amd64 $(ENTRYPOINT)
	GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/$(VERSION)/$(BINARY)_$(VERSION)_linux_arm64 $(ENTRYPOINT)
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/$(VERSION)/$(BINARY)_$(VERSION)_darwin_amd64 $(ENTRYPOINT)
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/$(VERSION)/$(BINARY)_$(VERSION)_darwin_arm64 $(ENTRYPOINT)
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/$(VERSION)/$(BINARY)_$(VERSION)_windows_amd64.exe $(ENTRYPOINT)
	@cd $(DIST_DIR)/$(VERSION) && sha256sum $(BINARY)_$(VERSION)_* > $(BINARY)_$(VERSION)_SHA256SUMS

test: test-unit test-e2e

test-unit:
	go test -race -v $$(go list ./... | grep -v /test/)

test-e2e: build
	go test -v ./test/e2e/...

test-e2e-coverage:
	@rm -rf $(COVERAGE_DIR)
	@mkdir -p $(COVERAGE_DIR)/raw
	go build $(GOFLAGS) -cover -ldflags "$(LDFLAGS)" -o $(BINARY) $(ENTRYPOINT)
	GOCOVERDIR=$(shell pwd)/$(COVERAGE_DIR)/raw go test -v ./test/e2e/...
	go tool covdata textfmt -i=$(COVERAGE_DIR)/raw -o=$(COVERAGE_DIR)/e2e.out
	go tool cover -html=$(COVERAGE_DIR)/e2e.out -o $(COVERAGE_DIR)/e2e.html

clean:
	rm -f $(BINARY)
	rm -rf $(DIST_DIR)
	rm -rf $(COVERAGE_DIR)

lint:
	go vet ./...
	staticcheck ./...
	gosec -quiet ./...

fmt:
	goimports -l -w .

install-tools:
	go install golang.org/x/tools/cmd/goimports@v0.42.0
	go install honnef.co/go/tools/cmd/staticcheck@v0.7.0
	go install github.com/securego/gosec/v2/cmd/gosec@v2.23.0

doc:
	@mkdir -p docs
	go run ./internal/tools/doc > docs/rules.md
