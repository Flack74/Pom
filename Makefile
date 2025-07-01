VERSION := $(shell git describe --tags --always --dirty)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse HEAD)

LDFLAGS := -ldflags "-X github.com/Flack74/pom/cmd.version=$(VERSION) -X github.com/Flack74/pom/cmd.buildDate=$(BUILD_DATE) -X github.com/Flack74/pom/cmd.gitCommit=$(GIT_COMMIT)"

.PHONY: build
build:
	@echo "Building pom with version info..."
	cd web/frontend && npm install && npm run build
	go build $(LDFLAGS) -o pom .

.PHONY: build-cli
build-cli:
	@echo "Building CLI only..."
	go build $(LDFLAGS) -o pom .

.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Build Date: $(BUILD_DATE)"
	@echo "Git Commit: $(GIT_COMMIT)"

.PHONY: clean
clean:
	rm -f pom
	rm -rf web/build/

.PHONY: install
install: build
	sudo cp pom /usr/local/bin/

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build      - Build with web UI and version info"
	@echo "  build-cli  - Build CLI only with version info"
	@echo "  version    - Show version info"
	@echo "  clean      - Clean build artifacts"
	@echo "  install    - Install to /usr/local/bin"
	@echo "  help       - Show this help"