.PHONY: all test clean compress

CGO_ENABLED ?= 0
GOFLAGS ?= -trimpath
LDFLAGS ?= -s -w

all: RUN INSTALL

RUN:
	CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o RUN ./cmd/run

INSTALL:
	CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o INSTALL ./cmd/install

test:
	go test ./internal/...

compress:
	upx --best --lzma RUN INSTALL || true

clean:
	rm -f RUN INSTALL
