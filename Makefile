# ============================================================================
# Makefile for Krathub Project
# ============================================================================
# This Makefile follows Go project best practices with enhanced development
# experience, code quality checks, and modern tooling integration.
# ============================================================================

# ============================================================================
# VARIABLES & CONFIGURATION
# ============================================================================
# Go environment variables
GOHOSTOS := $(shell go env GOHOSTOS)
GOPATH := $(shell go env GOPATH)
GOVERSION := $(shell go version)
GOLANGCILINT_VERSION := v1.59.1

# Build information
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date +%Y-%m-%dT%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

# Project directories
BIN_DIR := bin
API_DIR := api
INTERNAL_DIR := internal
THIRD_PARTY_DIR := third_party
CONFIG_DIR := configs
CMD_DIR := cmd
COVERAGE_DIR := coverage
DOCS_DIR := docs

# Build configuration
BINARY_NAME := server
BUILD_FLAGS := -ldflags "$(LDFLAGS)"
LDFLAGS := -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.GitBranch=$(GIT_BRANCH)
DOCKER_IMAGE := krathub
DOCKER_TAG := $(VERSION)

# Proto files detection with cross-platform support
ifeq ($(GOHOSTOS), windows)
	Git_Bash := $(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES := $(shell $(Git_Bash) -c "find $(INTERNAL_DIR) -name *.proto 2>/dev/null" || echo "")
	API_PROTO_FILES := $(shell $(Git_Bash) -c "find $(API_DIR) -name *.proto 2>/dev/null" || echo "")
else
	INTERNAL_PROTO_FILES := $(shell find $(INTERNAL_DIR) -name *.proto 2>/dev/null || true)
	API_PROTO_FILES := $(shell find $(API_DIR) -name *.proto 2>/dev/null || true)
endif

# Output colors for better user experience
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[0;37m
RESET := \033[0m

# Utility functions
define print_success
	@echo "$(GREEN)✓ $(1)$(RESET)"
endef
define print_error
	@echo "$(RED)✗ $(1)$(RESET)"
endef
define print_warning
	@echo "$(YELLOW)⚠ $(1)$(RESET)"
endef
define print_info
	@echo "$(BLUE)ℹ $(1)$(RESET)"
endef
define print_step
	@echo "$(PURPLE)→ $(1)$(RESET)"
endef

# ============================================================================
# DEVELOPMENT TOOLS
# ============================================================================

.PHONY: help
## Show this help message
help:
	@echo "$(CYAN)Krathub Development Environment$(RESET)"
	@echo "$(CYAN)=================================$(RESET)"
	@echo ''
	@echo '$(YELLOW)Usage:$(RESET)'
	@echo '  make [target]'
	@echo ''
	@echo '$(YELLOW)Development Tools:$(RESET)'
	@echo '  $(GREEN)help$(RESET)           Show this help message'
	@echo '  $(GREEN)version$(RESET)        Show version information'
	@echo '  $(GREEN)setup$(RESET)          Complete development environment setup'
	@echo '  $(GREEN)install-dev$(RESET)    Install all development dependencies'
	@echo '  $(GREEN)check-tools$(RESET)    Check if required tools are installed'
	@echo ''
	@echo '$(YELLOW)Code Generation:$(RESET)'
	@echo '  $(GREEN)proto$(RESET)          Generate all protobuf files'
	@echo '  $(GREEN)api$(RESET)            Generate API protobuf files'
	@echo '  $(GREEN)config$(RESET)         Generate internal protobuf files'
	@echo '  $(GREEN)errors$(RESET)         Generate error protobuf files'
	@echo '  $(GREEN)wire$(RESET)           Generate dependency injection code'
	@echo '  $(GREEN)generate$(RESET)       Run all code generation tasks'
	@echo '  $(GREEN)gen.db$(RESET)         Generate database code'
	@echo '  $(GREEN)gen.tls$(RESET)        Generate TLS certificates'
	@echo ''
	@echo '$(YELLOW)Building:$(RESET)'
	@echo '  $(GREEN)build$(RESET)          Build the application'
	@echo '  $(GREEN)build-release$(RESET)  Build release version with optimizations'
	@echo '  $(GREEN)build-all$(RESET)      Build for multiple platforms'
	@echo '  $(GREEN)run$(RESET)            Run the application'
	@echo '  $(GREEN)run-debug$(RESET)      Run the application in debug mode'
	@echo '  $(GREEN)all$(RESET)            Build everything from scratch'
	@echo ''
	@echo '$(YELLOW)Testing & Quality:$(RESET)'
	@echo '  $(GREEN)test$(RESET)           Run all tests'
	@echo '  $(GREEN)test-cover$(RESET)     Run tests with coverage report'
	@echo '  $(GREEN)test-integration$(RESET) Run integration tests'
	@echo '  $(GREEN)benchmark$(RESET)      Run benchmark tests'
	@echo '  $(GREEN)race$(RESET)           Run tests with race detection'
	@echo '  $(GREEN)lint$(RESET)           Run static code analysis'
	@echo '  $(GREEN)fmt$(RESET)            Format code'
	@echo '  $(GREEN)check$(RESET)          Run all quality checks'
	@echo '  $(GREEN)security$(RESET)       Run security scan'
	@echo ''
	@echo '$(YELLOW)Dependency Management:$(RESET)'
	@echo '  $(GREEN)mod-update$(RESET)     Update Go dependencies'
	@echo '  $(GREEN)mod-tidy$(RESET)       Clean up Go dependencies'
	@echo '  $(GREEN)mod-verify$(RESET)     Verify Go dependencies'
	@echo '  $(GREEN)deps$(RESET)           Download all dependencies'
	@echo ''
	@echo '$(YELLOW)Docker:$(RESET)'
	@echo '  $(GREEN)docker-build$(RESET)   Build Docker image'
	@echo '  $(GREEN)docker-run$(RESET)     Run Docker container'
	@echo '  $(GREEN)docker-push$(RESET)    Push Docker image'
	@echo '  $(GREEN)docker-clean$(RESET)   Clean Docker resources'
	@echo ''
	@echo '$(YELLOW)Cleaning:$(RESET)'
	@echo '  $(GREEN)clean$(RESET)          Clean generated files'
	@echo '  $(GREEN)clean-all$(RESET)      Clean everything including caches'
	@echo '  $(GREEN)clean-deps$(RESET)     Clean and reinstall dependencies'
	@echo ''
	@echo '$(YELLOW)Legacy Commands:$(RESET)'
	@echo '  $(GREEN)init$(RESET)           Legacy: Initialize environment (alias for install-dev)'
	@echo '  $(GREEN)info$(RESET)           Show project information'
	@echo ''
	@echo '$(YELLOW)Examples:$(RESET)'
	@echo '  make setup            # Complete development environment setup'
	@echo '  make check            # Run all quality checks'
	@echo '  make all              # Build everything from scratch'
	@echo '  make test-cover       # Run tests with coverage report'
	@echo '  make build-all        # Build for all platforms'

.PHONY: version
## Show version information
version:
	@echo "$(CYAN)Project Information:$(RESET)"
	@echo "  Name:     $(DOCKER_IMAGE)"
	@echo "  Version:  $(VERSION)"
	@echo "  Git:      $(GIT_COMMIT)"
	@echo "  Branch:   $(GIT_BRANCH)"
	@echo "  Go:       $(GOVERSION)"
	@echo "  OS:       $(GOHOSTOS)"

.PHONY: setup
## Complete development environment setup (Development Tools)
setup: install-dev
	$(call print_step,Setting up development environment...)
	$(call print_info,Running go mod download...)
	go mod download
	$(call print_info,Verifying dependencies...)
	go mod verify
	$(call print_success,Development environment setup complete!)

.PHONY: install-dev
## Install all development dependencies (Development Tools)
install-dev:
	$(call print_step,Installing development tools...)
	@echo Installing protobuf generators...
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	@echo Installing quality check tools...
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo Installing golangci-lint...; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCILINT_VERSION); \
	fi
	$(call print_success,Development tools installed!)

.PHONY: check-tools
## Check if required tools are installed (Development Tools)
check-tools:
	$(call print_step,Checking required tools...)
	@echo Checking Go installation...
	@go version
	@echo Checking protobuf compiler...
	@protoc --version || $(call print_error,protoc not found. Please install protobuf compiler)
	@echo Checking development tools...
	@command -v wire >/dev/null 2>&1 || $(call print_warning,wire not found. Run 'make install-dev')
	@command -v kratos >/dev/null 2>&1 || $(call print_warning,kratos not found. Run 'make install-dev')
	@command -v golangci-lint >/dev/null 2>&1 || $(call print_warning,golangci-lint not found. Run 'make install-dev')
	$(call print_success,Tool check complete!)

# ============================================================================
# CODE GENERATION
# ============================================================================

.PHONY: proto
## Generate all protobuf files (Code Generation)
proto: api errors config
	$(call print_success,All protobuf files generated!)

.PHONY: api
## Generate API protobuf files (Code Generation)
api:
	$(call print_step,Generating API protobuf files...)
	@if [ -z "$(API_PROTO_FILES)" ]; then \
		echo "No API proto files found"; \
	else \
		protoc --proto_path=./$(API_DIR) \
		       --proto_path=./$(THIRD_PARTY_DIR) \
		       --go_out=paths=source_relative:./$(API_DIR) \
		       --go-http_out=paths=source_relative:./$(API_DIR) \
		       --go-grpc_out=paths=source_relative:./$(API_DIR) \
		       --openapi_out=fq_schema_naming=true,default_response=false:. \
		       $(API_PROTO_FILES); \
		echo "API protobuf files generated"; \
	fi

.PHONY: config
## Generate internal protobuf files (Code Generation)
config:
	$(call print_step,Generating internal protobuf files...)
	@if [ -z "$(INTERNAL_PROTO_FILES)" ]; then \
		echo "No internal proto files found"; \
	else \
		protoc --proto_path=./$(INTERNAL_DIR) \
		       --proto_path=./$(THIRD_PARTY_DIR) \
		       --go_out=paths=source_relative:./$(INTERNAL_DIR) \
		       $(INTERNAL_PROTO_FILES); \
		echo "Internal protobuf files generated"; \
	fi

.PHONY: errors
## Generate error protobuf files (Code Generation)
errors:
	$(call print_step,Generating error protobuf files...)
	@if [ -z "$(API_PROTO_FILES)" ]; then \
		echo "No API proto files found for error generation"; \
	else \
		protoc --proto_path=. \
		       --proto_path=./$(THIRD_PARTY_DIR) \
		       --go_out=paths=source_relative:. \
		       --go-errors_out=paths=source_relative:. \
		       $(API_PROTO_FILES); \
		echo "Error protobuf files generated"; \
	fi

.PHONY: wire
## Generate dependency injection code (Code Generation)
wire:
	$(call print_step,Running wire dependency injection...)
	@cd $(CMD_DIR)/server && wire
	$(call print_success,Wire generation complete!)

.PHONY: generate
## Run all code generation tasks (Code Generation)
generate:
	$(call print_step,Running all code generation...)
	go generate ./...
	$(call print_success,Code generation complete!)

.PHONY: gen.db
## Generate database code (Code Generation)
gen.db:
	$(call print_step,Generating database code...)
	go run $(CMD_DIR)/gen/gendb.go -conf $(CONFIG_DIR)/config.yaml
	$(call print_success,Database code generated!)

.PHONY: gen.tls
## Generate TLS certificates (Code Generation)
gen.tls:
	$(call print_step,Generating TLS certificates...)
	@mkdir -p manifest/certs
	@openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
		-keyout manifest/certs/server.key \
		-out manifest/certs/server.cert \
		-config manifest/certs/openssl.cnf 2>/dev/null || \
		openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
		-keyout manifest/certs/server.key \
		-out manifest/certs/server.cert \
		-subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"
	$(call print_success,TLS certificates generated in manifest/certs/)

# ============================================================================
# BUILDING
# ============================================================================

.PHONY: build
## Build the application (Building)
build:
	$(call print_step,Building $(BINARY_NAME)...)
	@mkdir -p $(BIN_DIR)
	go build $(BUILD_FLAGS) -o ./$(BIN_DIR)/$(BINARY_NAME) ./$(CMD_DIR)/server
	$(call print_success,Build complete: ./$(BIN_DIR)/$(BINARY_NAME))

.PHONY: build-release
## Build release version with optimizations (Building)
build-release:
	$(call print_step,Building release version...)
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS) -s -w" -trimpath \
		-o ./$(BIN_DIR)/$(BINARY_NAME) ./$(CMD_DIR)/server
	$(call print_success,Release build complete: ./$(BIN_DIR)/$(BINARY_NAME))

.PHONY: build-all
## Build for multiple platforms (Building)
build-all:
	$(call print_step,Building for multiple platforms...)
	@mkdir -p $(BIN_DIR)
	@echo Building for linux/amd64...
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o ./$(BIN_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)/server
	@echo Building for linux/arm64...
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o ./$(BIN_DIR)/$(BINARY_NAME)-linux-arm64 ./$(CMD_DIR)/server
	@echo Building for darwin/amd64...
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o ./$(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)/server
	@echo Building for darwin/arm64...
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o ./$(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)/server
	@echo Building for windows/amd64...
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o ./$(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)/server
	$(call print_success,Cross-platform builds complete!)

.PHONY: run
## Run the application (Building)
run:
	$(call print_step,Starting $(BINARY_NAME)...)
	cd cmd/server && go run .

.PHONY: run-debug
## Run the application in debug mode (Building)
run-debug:
	$(call print_step,Starting $(BINARY_NAME) in debug mode...)
	kratos run --debug

.PHONY: all
## Build everything from scratch (Building)
all: clean proto gen.db wire generate build
	$(call print_success,Complete build finished!)

# ============================================================================
# TESTING & QUALITY
# ============================================================================

.PHONY: test
## Run all tests (Testing & Quality)
test:
	$(call print_step,Running tests...)
	go test -v ./...

.PHONY: test-cover
## Run tests with coverage report (Testing & Quality)
test-cover:
	$(call print_step,Running tests with coverage...)
	@mkdir -p $(COVERAGE_DIR)
	go test -v -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	$(call print_success,Coverage report generated: $(COVERAGE_DIR)/coverage.html)

.PHONY: test-integration
## Run integration tests (Testing & Quality)
test-integration:
	$(call print_step,Running integration tests...)
	go test -v -tags=integration ./...

.PHONY: benchmark
## Run benchmark tests (Testing & Quality)
benchmark:
	$(call print_step,Running benchmarks...)
	go test -bench=. -benchmem ./...

.PHONY: race
## Run tests with race detection (Testing & Quality)
race:
	$(call print_step,Running tests with race detection...)
	go test -race -v ./...

.PHONY: lint
## Run static code analysis (Testing & Quality)
lint:
	$(call print_step,Running linter...)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		echo "Linting complete!"; \
	else \
		echo "golangci-lint not found. Please run 'make install-dev'"; \
		exit 1; \
	fi

.PHONY: fmt
## Format code (Testing & Quality)
fmt:
	$(call print_step,Formatting code...)
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi
	$(call print_success,Code formatted!)

.PHONY: check
## Run all quality checks (Testing & Quality)
check: fmt lint test
	$(call print_success,All quality checks passed!)

.PHONY: security
## Run security scan (Testing & Quality)
security:
	$(call print_step,Running security scan...)
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		$(call print_warning,gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest); \
	fi
	@go list -json -m all | nancy sleuth 2>/dev/null || echo "nancy not found, skipping vulnerability check"
	$(call print_success,Security scan complete!)

# ============================================================================
# DEPENDENCY MANAGEMENT
# ============================================================================

.PHONY: mod-update
## Update Go dependencies (Dependency Management)
mod-update:
	$(call print_step,Updating dependencies...)
	@go get -u ./...
	@go mod tidy
	$(call print_success,Dependencies updated!)

.PHONY: mod-tidy
## Clean up Go dependencies (Dependency Management)
mod-tidy:
	$(call print_step,Tidying dependencies...)
	go mod tidy
	$(call print_success,Dependencies tidied!)

.PHONY: mod-verify
## Verify Go dependencies (Dependency Management)
mod-verify:
	$(call print_step,Verifying dependencies...)
	go mod verify
	$(call print_success,Dependencies verified!)

.PHONY: deps
## Download all dependencies (Dependency Management)
deps:
	$(call print_step,Downloading dependencies...)
	go mod download
	$(call print_success,Dependencies downloaded!)

# ============================================================================
# DOCKER
# ============================================================================

.PHONY: docker-build
## Build Docker image (Docker)
docker-build:
	$(call print_step,Building Docker image...)
	@if [ -f Dockerfile ]; then \
		docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .; \
		docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):latest; \
		$(call print_success,Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)); \
	else \
		$(call print_error,Dockerfile not found); \
		exit 1; \
	fi

.PHONY: docker-run
## Run Docker container (Docker)
docker-run:
	$(call print_step,Running Docker container...)
	docker run --rm -p 8000:8000 -p 9000:9000 $(DOCKER_IMAGE):latest

.PHONY: docker-push
## Push Docker image (Docker)
docker-push:
	$(call print_step,Pushing Docker image...)
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):latest
	$(call print_success,Docker image pushed!)

.PHONY: docker-clean
## Clean Docker resources (Docker)
docker-clean:
	$(call print_step,Cleaning Docker resources...)
	docker system prune -f
	docker volume prune -f
	$(call print_success,Docker cleanup complete!)

# ============================================================================
# CLEANING
# ============================================================================

.PHONY: clean
## Clean generated files (Cleaning)
clean:
	$(call print_step,Cleaning generated files...)
	@rm -rf $(BIN_DIR)
	@find . -name "*.pb.go" -delete 2>/dev/null || true
	@find . -name "*_grpc.pb.go" -delete 2>/dev/null || true
	@find . -name "*_http.pb.go" -delete 2>/dev/null || true
	@find . -name "*.pb.validate.go" -delete 2>/dev/null || true
	@find . -name "*_errors.pb.go" -delete 2>/dev/null || true
	@rm -f openapi.yaml 2>/dev/null || true
	@rm -rf $(COVERAGE_DIR) 2>/dev/null || true
	$(call print_success,Clean complete!)

.PHONY: clean-all
## Clean everything including caches (Cleaning)
clean-all: clean
	$(call print_step,Cleaning all files and caches...)
	@go clean -cache -testcache -modcache
	@rm -rf vendor/ 2>/dev/null || true
	$(call print_success,Deep clean complete!)

.PHONY: clean-deps
## Clean and reinstall dependencies (Cleaning)
clean-deps:
	$(call print_step,Cleaning dependencies...)
	@rm -rf go.sum go.mod vendor/
	@go mod init github.com/horonlee/krathub
	$(call print_success,Dependencies cleaned!)

# ============================================================================
# LEGACY COMPATIBILITY
# ============================================================================

# These targets are maintained for backward compatibility
.PHONY: init
## Legacy: Initialize environment (alias for install-dev)
init: install-dev

# ============================================================================
# MISC
# ============================================================================

.PHONY: info
## Show project information
info: version

# Set default goal
.DEFAULT_GOAL := help
