# ============================================================================
# Makefile for Krathub Project
# ============================================================================
# Based on go-wind-admin project structure
# ============================================================================

ifeq ($(OS),Windows_NT)
    IS_WINDOWS := 1
endif

# load environment variables from .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# ============================================================================
# VARIABLES & CONFIGURATION
# ============================================================================

# Directories
CURRENT_DIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
ROOT_DIR    := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
API_DIR     := api
PKG_DIR     := pkg

# Find all service Makefiles in app directory
SRCS_MK := $(foreach dir, app, $(wildcard $(dir)/*/*/Makefile))

# Go environment
GOPATH := $(shell go env GOPATH)
GOVERSION := $(shell go version)

# Build information
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date +%Y-%m-%dT%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")

# LDFLAGS
LDFLAGS := -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.GitBranch=$(GIT_BRANCH)

# Output colors
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
CYAN := \033[0;36m
RESET := \033[0m

# ============================================================================
# MAIN TARGETS
# ============================================================================

.PHONY: help env init plugin cli dep vendor test cover vet lint
.PHONY: wire ent gen api openapi build build_only all docker clean

# show environment variables
env:
	@echo "CURRENT_DIR: $(CURRENT_DIR)"
	@echo "ROOT_DIR: $(ROOT_DIR)"
	@echo "SRCS_MK: $(SRCS_MK)"
	@echo "VERSION: $(VERSION)"
	@echo "GOVERSION: $(GOVERSION)"

# initialize develop environment
init: plugin cli
	@echo "$(GREEN)✓ Development environment initialized$(RESET)"

# install protoc plugins
plugin:
	@echo "$(CYAN)Installing protoc plugins...$(RESET)"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	@go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	@go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	@go install github.com/envoyproxy/protoc-gen-validate@latest
	@echo "$(GREEN)✓ Protoc plugins installed$(RESET)"

# install cli tools
cli:
	@echo "$(CYAN)Installing CLI tools...$(RESET)"
	@go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	@go install github.com/google/gnostic@latest
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/google/wire/cmd/wire@latest
	@echo "$(GREEN)✓ CLI tools installed$(RESET)"

# download dependencies of module
dep:
	@go mod download

# create vendor
vendor:
	@go mod vendor

# run tests
test:
	@go test ./...

# run coverage tests
cover:
	@go test -v ./... -coverprofile=coverage.out

# run static analysis
vet:
	@go vet ./...

# run lint
lint:
	@golangci-lint run

# generate wire code for all services
wire:
	@echo "$(CYAN)Generating wire code for all services...$(RESET)"
	@$(foreach dir, $(dir $(realpath $(SRCS_MK))),\
      cd $(dir) && make wire;\
    )
	@echo "$(GREEN)✓ Wire code generated$(RESET)"

# generate all code
gen: wire api openapi
	@echo "$(GREEN)✓ All code generated$(RESET)"

# generate protobuf api go code
api:
	@echo "$(CYAN)Generating protobuf Go code...$(RESET)"
	@cd $(API_DIR) && buf generate
	@echo "$(GREEN)✓ Protobuf Go code generated$(RESET)"

# generate protobuf api OpenAPI v3 docs for all services
openapi:
	@echo "$(CYAN)Generating OpenAPI documentation for all services...$(RESET)"
	@$(foreach dir, $(dir $(realpath $(SRCS_MK))),\
      cd $(dir) && make openapi;\
    )
	@echo "$(GREEN)✓ OpenAPI documentation generated$(RESET)"

# lint protobuf files
lint.proto:
	@echo "$(CYAN)Linting protobuf files...$(RESET)"
	@cd $(API_DIR) && buf lint
	@echo "$(GREEN)✓ Proto lint complete$(RESET)"

# update buf dependencies
buf.update:
	@echo "$(CYAN)Updating buf dependencies...$(RESET)"
	@cd $(API_DIR)/protos && buf dep update
	@echo "$(GREEN)✓ Buf dependencies updated$(RESET)"

# build all service applications
build: api openapi
	@echo "$(CYAN)Building all services...$(RESET)"
	@$(foreach dir, $(dir $(realpath $(SRCS_MK))),\
      cd $(dir) && make build;\
    )
	@echo "$(GREEN)✓ All services built$(RESET)"

# only build all service applications without generating api and openapi
build_only:
	@echo "$(CYAN)Building all services (without code generation)...$(RESET)"
	@$(foreach dir, $(dir $(realpath $(SRCS_MK))),\
      cd $(dir) && make build_only;\
    )
	@echo "$(GREEN)✓ All services built$(RESET)"

# generate & build all service applications
all:
	@echo "$(CYAN)Generating and building all services...$(RESET)"
	@$(foreach dir, $(dir $(realpath $(SRCS_MK))),\
      cd $(dir) && make app;\
    )
	@echo "$(GREEN)✓ All services generated and built$(RESET)"

# build docker images for all services
docker:
	@echo "$(CYAN)Building Docker images...$(RESET)"
	@$(foreach dir, $(dir $(realpath $(SRCS_MK))),\
      cd $(dir) && make docker;\
    )
	@echo "$(GREEN)✓ Docker images built$(RESET)"

# clean build artifacts
clean:
	@echo "$(CYAN)Cleaning build artifacts...$(RESET)"
	@rm -rf $(API_DIR)/gen
	@$(foreach dir, $(dir $(realpath $(SRCS_MK))),\
      cd $(dir) && make clean;\
    )
	@echo "$(GREEN)✓ Clean complete$(RESET)"

# show help
help:
	@echo ""
	@echo "$(CYAN)Krathub Development Environment$(RESET)"
	@echo "$(CYAN)=================================$(RESET)"
	@echo ""
	@echo "Usage:"
	@echo " make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "  $(GREEN)%-15s$(RESET) %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
