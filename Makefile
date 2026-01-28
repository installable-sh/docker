.PHONY: all test clean compress fmt lint check setup

CGO_ENABLED ?= 0
GOFLAGS ?= -trimpath
LDFLAGS ?= -s -w

all: setup fmt lint test RUN INSTALL

RUN:
	@echo "==> Building RUN..."
	@CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o RUN ./cmd/run

INSTALL:
	@echo "==> Building INSTALL..."
	@CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o INSTALL ./cmd/install

test:
	@echo "==> Running tests..."
	@go test ./internal/...

fmt:
	@echo "==> Formatting code..."
	@gofmt -s -w .

lint:
	@echo "==> Running linter..."
	@golangci-lint run

check:
	@echo "==> Checking formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Run 'make fmt' to fix formatting" && gofmt -l . && exit 1)

compress:
	@echo "==> Compressing binaries..."
	@upx --best --lzma RUN INSTALL 2>/dev/null || true

clean:
	@echo "==> Cleaning..."
	@rm -f RUN INSTALL
	@go clean -testcache

setup:
	@echo "==> Setting up development environment..."
	@[ -f scripts/pre-commit ] && [ -d .git/hooks ] && ln -sf ../../scripts/pre-commit .git/hooks/pre-commit || true
	@command -v golangci-lint >/dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
