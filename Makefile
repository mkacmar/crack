.PHONY: build build-release test clean lint fmt install-tools

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS = -X github.com/mkacmar/crack/internal/version.Version=$(VERSION) \
          -X github.com/mkacmar/crack/internal/version.GitCommit=$(COMMIT) \
          -X github.com/mkacmar/crack/internal/version.BuildTime=$(BUILD_TIME)
RELEASE_LDFLAGS = -s -w $(LDFLAGS)
DIST_DIR = dist

build:
	go build -ldflags "$(LDFLAGS)" -o crack ./cmd/crack

build-release:
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/crack-linux-amd64 ./cmd/crack
	GOOS=linux GOARCH=arm64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/crack-linux-arm64 ./cmd/crack
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/crack-darwin-amd64 ./cmd/crack
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/crack-darwin-arm64 ./cmd/crack
	GOOS=windows GOARCH=amd64 go build -ldflags "$(RELEASE_LDFLAGS)" -o $(DIST_DIR)/crack-windows-amd64.exe ./cmd/crack
	@echo ""
	@echo "SHA256 Checksums:"
	@cd $(DIST_DIR) && sha256sum *

test:
	go test -v ./...

clean:
	rm -f crack
	rm -rf $(DIST_DIR)

lint:
	go vet ./...
	staticcheck ./...

fmt:
	goimports -l -w .

install-tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest