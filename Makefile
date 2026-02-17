VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BINARY  := odds
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build test lint install clean

build:
	go build $(LDFLAGS) -o $(BINARY) .

test:
	go test ./...

lint:
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping"; \
	fi

install:
	go install $(LDFLAGS) .

clean:
	rm -f $(BINARY)
