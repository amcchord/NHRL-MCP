# NHRL MCP Server Makefile

BINARY_NAME=nhrl-mcp-server
VERSION=v1.8.1
BUILD_DIR=build

# Code signing variables (set these via environment or command line)
DEVELOPER_ID ?= ""
KEYCHAIN_PROFILE ?= ""

# Build flags for optimized binaries
LDFLAGS=-ldflags "-s -w"

# Default target
.PHONY: all
all: build

# Build for current platform
.PHONY: build
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Build for all platforms
.PHONY: build-all
build-all: clean
	mkdir -p $(BUILD_DIR)
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	# Windows ARM64
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe .

# Build for all platforms with macOS code signing
.PHONY: build-all-signed
build-all-signed: clean
	@if [ -z "$(DEVELOPER_ID)" ]; then \
		echo "Error: DEVELOPER_ID must be set for code signing"; \
		echo "Usage: make build-all-signed DEVELOPER_ID='Developer ID Application: Your Name (TEAMID)'"; \
		exit 1; \
	fi
	mkdir -p $(BUILD_DIR)
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	$(MAKE) sign-macos BINARY=$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	$(MAKE) sign-macos BINARY=$(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	# Windows ARM64
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe .

# Sign macOS binary
.PHONY: sign-macos
sign-macos:
	@if [ -z "$(BINARY)" ]; then \
		echo "Error: BINARY must be specified"; \
		exit 1; \
	fi
	@if [ -z "$(DEVELOPER_ID)" ]; then \
		echo "Error: DEVELOPER_ID must be set"; \
		exit 1; \
	fi
	@echo "Signing $(BINARY) with identity: $(DEVELOPER_ID)"
	codesign --sign "$(DEVELOPER_ID)" \
		--options runtime \
		--timestamp \
		--verbose=2 \
		"$(BINARY)"
	@echo "Verifying signature for $(BINARY)"
	codesign --verify --verbose=2 "$(BINARY)"
	@echo "Note: spctl assessment may fail for unnotarized binaries"
	-spctl --assess --type execute --verbose=2 "$(BINARY)"

# Notarize macOS binaries (requires Apple Developer account)
.PHONY: notarize-macos
notarize-macos:
	@if [ -z "$(KEYCHAIN_PROFILE)" ]; then \
		echo "Error: KEYCHAIN_PROFILE must be set for notarization"; \
		echo "Create a keychain profile with: xcrun notarytool store-credentials"; \
		exit 1; \
	fi
	@echo "Creating ZIP archives for notarization..."
	cd $(BUILD_DIR) && \
	zip $(BINARY_NAME)-darwin-amd64.zip $(BINARY_NAME)-darwin-amd64 && \
	zip $(BINARY_NAME)-darwin-arm64.zip $(BINARY_NAME)-darwin-arm64
	@echo "Submitting for notarization..."
	xcrun notarytool submit $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64.zip \
		--keychain-profile "$(KEYCHAIN_PROFILE)" \
		--wait
	xcrun notarytool submit $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64.zip \
		--keychain-profile "$(KEYCHAIN_PROFILE)" \
		--wait
	@echo "Stapling notarization..."
	@echo "Note: Stapling may fail for some binary types, but notarization is still valid"
	-xcrun stapler staple $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
	-xcrun stapler staple $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64

# Install dependencies
.PHONY: deps
deps:
	go mod tidy
	go mod download

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run the server with test credentials (for development)
.PHONY: test-run
test-run:
	@if [ -z "$(NHRL_API_USER_ID)" ] || [ -z "$(NHRL_API_KEY)" ]; then \
		echo "Error: Please set NHRL_API_USER_ID and NHRL_API_KEY environment variables"; \
		echo "Or source your environment file: source scripts/dev-env.sh"; \
		exit 1; \
	fi
	go run . -exit-after-first -tools full

# Test tools functionality
.PHONY: test-tools
test-tools:
	@echo "Testing tools registration..."
	@if [ -z "$(NHRL_API_USER_ID)" ] || [ -z "$(NHRL_API_KEY)" ]; then \
		echo "Warning: API credentials not set, using dummy values for tool testing"; \
		NHRL_API_USER_ID=test NHRL_API_KEY=test go run . -exit-after-first -tools full 2>&1 | grep -E "(tool|Tool)" || true; \
	else \
		go run . -exit-after-first -tools full 2>&1 | grep -E "(tool|Tool)" || true; \
	fi

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Run the server locally
.PHONY: run
run:
	go run .

# Install to system PATH
.PHONY: install
install: build
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Install for testing (current platform only)
.PHONY: install-test
install-test:
	./scripts/install_for_test.sh

# Create release packages
.PHONY: release
release: build-all
	cd $(BUILD_DIR) && \
	tar -czf $(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 && \
	tar -czf $(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64 && \
	tar -czf $(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 && \
	tar -czf $(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64 && \
	zip $(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe && \
	zip $(BINARY_NAME)-$(VERSION)-windows-arm64.zip $(BINARY_NAME)-windows-arm64.exe

# Create signed release packages (macOS binaries will be signed and notarized)
.PHONY: release-signed
release-signed: build-all-signed notarize-macos
	cd $(BUILD_DIR) && \
	tar -czf $(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 && \
	tar -czf $(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64 && \
	tar -czf $(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 && \
	tar -czf $(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64 && \
	zip $(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe && \
	zip $(BINARY_NAME)-$(VERSION)-windows-arm64.zip $(BINARY_NAME)-windows-arm64.exe

# Generate checksums for release
.PHONY: checksums
checksums:
	@if [ ! -d "$(BUILD_DIR)" ]; then \
		echo "Error: Build directory does not exist. Run 'make build-all' first."; \
		exit 1; \
	fi
	cd $(BUILD_DIR) && \
	find . -name "$(BINARY_NAME)-$(VERSION)-*" -type f -exec sha256sum {} \; > checksums.sha256 && \
	echo "Checksums saved to $(BUILD_DIR)/checksums.sha256"

# Verify checksums
.PHONY: verify-checksums
verify-checksums:
	@if [ ! -f "$(BUILD_DIR)/checksums.sha256" ]; then \
		echo "Error: Checksums file not found. Run 'make checksums' first."; \
		exit 1; \
	fi
	cd $(BUILD_DIR) && sha256sum -c checksums.sha256

# Show build info
.PHONY: info
info:
	@echo "NHRL MCP Server Build Information"
	@echo "=================================="
	@echo "Binary Name:     $(BINARY_NAME)"
	@echo "Version:         $(VERSION)"
	@echo "Build Directory: $(BUILD_DIR)"
	@echo "Go Version:      $$(go version)"
	@echo ""
	@echo "Available targets for cross-compilation:"
	@echo "  - linux/amd64"
	@echo "  - linux/arm64"
	@echo "  - darwin/amd64 (Intel Mac)"
	@echo "  - darwin/arm64 (Apple Silicon)"
	@echo "  - windows/amd64"
	@echo "  - windows/arm64"

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		go vet ./...; \
	fi

# Run security scan
.PHONY: security
security:
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Full CI pipeline
.PHONY: ci
ci: deps fmt lint test build-all checksums

# Show help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build              - Build server for current platform"
	@echo "  build-all          - Build server for all supported platforms"
	@echo "  build-all-signed   - Build server for all platforms with macOS code signing"
	@echo "  sign-macos         - Sign a specific macOS binary"
	@echo "  notarize-macos     - Notarize macOS binaries with Apple"
	@echo "  deps               - Install/update dependencies"
	@echo "  test               - Run Go tests"
	@echo "  test-run           - Run server with test credentials (exits after first request)"
	@echo "  test-tools         - Test tools registration"
	@echo "  clean              - Remove all build artifacts"
	@echo "  run                - Run the server locally"
	@echo "  install            - Install to system PATH"
	@echo "  install-test       - Install for testing (current platform only)"
	@echo "  release            - Create release packages"
	@echo "  release-signed     - Create signed release packages"
	@echo "  checksums          - Generate SHA256 checksums for release packages"
	@echo "  verify-checksums   - Verify SHA256 checksums"
	@echo "  info               - Show build information"
	@echo "  fmt                - Format Go code"
	@echo "  lint               - Run code linter"
	@echo "  security           - Run security scan"
	@echo "  ci                 - Run full CI pipeline (deps, fmt, lint, test, build, checksums)"
	@echo "  help               - Show this help"
	@echo ""
	@echo "Code Signing Usage:"
	@echo "  make build-all-signed DEVELOPER_ID='Developer ID Application: Your Name (TEAMID)'"
	@echo "  make notarize-macos KEYCHAIN_PROFILE='YourKeychainProfile'"
	@echo ""
	@echo "Environment Variables:"
	@echo "  NHRL_API_USER_ID  - NHRL API User ID"
	@echo "  NHRL_API_KEY      - NHRL API Key"
	@echo "  DEVELOPER_ID            - macOS code signing identity"
	@echo "  KEYCHAIN_PROFILE        - macOS notarization keychain profile" 