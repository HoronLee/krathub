# Variables
GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S)
GIT_COMMIT=$(shell git rev-parse HEAD)

# Directories
BIN_DIR=bin
API_DIR=api
INTERNAL_DIR=internal
THIRD_PARTY_DIR=third_party
CONFIG_DIR=configs
CMD_DIR=cmd

# Proto files detection
ifeq ($(GOHOSTOS), windows)
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find $(INTERNAL_DIR) -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find $(API_DIR) -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find $(INTERNAL_DIR) -name *.proto 2>/dev/null)
	API_PROTO_FILES=$(shell find $(API_DIR) -name *.proto 2>/dev/null)
endif

# Build flags
LDFLAGS=-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./$(INTERNAL_DIR) \
	       --proto_path=./$(THIRD_PARTY_DIR) \
	       --go_out=paths=source_relative:./$(INTERNAL_DIR) \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./$(API_DIR) \
	       --proto_path=./$(THIRD_PARTY_DIR) \
	       --go_out=paths=source_relative:./$(API_DIR) \
	       --go-http_out=paths=source_relative:./$(API_DIR) \
	       --go-grpc_out=paths=source_relative:./$(API_DIR) \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)

.PHONY: errors
# generate errors proto
errors:
	protoc --proto_path=. \
	       --proto_path=./$(THIRD_PARTY_DIR) \
	       --go_out=paths=source_relative:. \
	       --go-errors_out=paths=source_relative:. \
	       $(API_PROTO_FILES)

.PHONY: proto
# generate all proto related files
proto: api errors config

.PHONY: build
# build
build:
	mkdir -p $(BIN_DIR) && go build -ldflags "$(LDFLAGS)" -o ./$(BIN_DIR)/server ./$(CMD_DIR)/server

.PHONY: wire
# generate wire
wire:
	cd $(CMD_DIR)/server && wire

.PHONY: generate
# generate
generate:
	go generate ./...
	go mod tidy

.PHONY: gen.db
# generate db
gen.db:
	go run $(CMD_DIR)/gen/gendb.go -conf $(CONFIG_DIR)/config.yaml

.PHONY: fmt
# format code
fmt:
	go fmt ./...

.PHONY: test
# run tests
test:
	go test -v ./...

.PHONY: clean
# clean generated files
clean:
	rm -rf $(BIN_DIR)
	find . -name "*.pb.go" -delete
	find . -name "*_grpc.pb.go" -delete
	find . -name "*_http.pb.go" -delete
	find . -name "*.pb.validate.go" -delete
	find . -name "*_errors.pb.go" -delete
	rm -f openapi.yaml

.PHONY: run
# run
run:
	kratos run

.PHONY: all
# generate all
all: clean proto gendb wire generate build

.PHONY: gen.tls
# generate tls cert and key
gen.tls:
	@mkdir -p manifest/certs
	@openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout manifest/certs/server.key -out manifest/certs/server.cert -config manifest/certs/openssl.cnf
	@echo "TLS certificate and key generated in manifest/certs/"

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($1, 0, index($1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "  %%-20s %%s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
