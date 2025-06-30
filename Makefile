.PHONY: all build install clean test lint package

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=pom
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X pom/cmd.version=$(VERSION) -X pom/cmd.buildDate=$(BUILD_TIME)"

# Installation paths
PREFIX=/usr/local
BINDIR=$(PREFIX)/bin
MANDIR=$(PREFIX)/share/man/man1

all: build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)

install: build
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 $(BINARY_NAME) $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	install -d $(DESTDIR)$(MANDIR)
	install -m 644 packaging/man/pom.1 $(DESTDIR)$(MANDIR)/pom.1
	gzip -f $(DESTDIR)$(MANDIR)/pom.1

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	rm -f $(DESTDIR)$(MANDIR)/pom.1.gz

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf dist/

test:
	$(GOTEST) -v ./...

lint:
	golangci-lint run

tidy:
	$(GOMOD) tidy

package:
	./scripts/build-packages.sh

# Development helpers
dev-deps:
	$(GOGET) golang.org/x/term
	$(GOGET) github.com/spf13/cobra
	golangci-lint --version >/dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin

help:
	@echo "Available targets:"
	@echo "  all        - Build the application (default)"
	@echo "  build      - Build the binary"
	@echo "  install    - Install the application to $(PREFIX)"
	@echo "  uninstall  - Remove the application from $(PREFIX)"
	@echo "  clean      - Remove build artifacts"
	@echo "  test       - Run tests"
	@echo "  lint       - Run linter"
	@echo "  tidy       - Tidy go.mod"
	@echo "  package    - Build all package formats"
	@echo "  dev-deps   - Install development dependencies"
	@echo "  help       - Show this help message" 