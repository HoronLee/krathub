# Krathub Development Guide

This project follows the [go-wind-admin](https://github.com/tx7do/go-wind-admin) project structure and best practices.

## Project Structure

```
krathub/
├── api/                                    # API definitions
│   ├── protos/                            # Protocol buffer definitions
│   │   ├── krathub/service/v1/           # Krathub HTTP services (i_*.proto)
│   │   ├── auth/service/v1/              # Auth gRPC service
│   │   ├── user/service/v1/              # User gRPC service
│   │   ├── test/service/v1/              # Test gRPC service
│   │   └── sayhello/service/v1/          # SayHello gRPC service
│   ├── gen/go/                           # Generated Go code
│   ├── buf.gen.yaml                      # Buf code generation config
│   └── buf.krathub.openapi.gen.yaml     # Krathub OpenAPI generation config
├── app/                                   # Application services
│   └── krathub/service/                  # Krathub service
│       ├── cmd/                          # Command line entry points
│       ├── internal/                     # Internal code
│       ├── configs/                      # Configuration files
│       ├── bin/                          # Build output
│       ├── openapi.yaml                  # Generated OpenAPI documentation
│       └── Makefile                      # Service Makefile (includes app.mk)
├── pkg/                                   # Shared packages
├── app.mk                                 # Common Makefile for all services
├── Makefile                               # Root Makefile
└── .env.example                          # Environment variables example

```

## Quick Start

### 1. Initialize Development Environment

```bash
# Install all required tools
make init

# Or install separately
make plugin  # Install protoc plugins
make cli     # Install CLI tools
```

### 2. Generate Code

```bash
# Generate all code (protobuf + OpenAPI)
make gen

# Or generate separately
make api      # Generate protobuf Go code
make openapi  # Generate OpenAPI documentation
make wire     # Generate dependency injection code
```

### 3. Build Services

```bash
# Build all services
make build

# Build without code generation
make build_only

# Generate and build everything
make all
```

### 4. Run Service

```bash
# Run krathub service
cd app/krathub/service
make run

# Or build and run
make build
./bin/server -c ./configs
```

## Development Workflow

### Adding a New Service

1. Create service directory structure:

```bash
mkdir -p app/{service-name}/service/{cmd/server,internal,configs}
```

1. Create service Makefile:

```bash
echo "include ../../../app.mk" > app/{service-name}/service/Makefile
```

1. Create OpenAPI config for the service:

```bash
cp api/buf.krathub.openapi.gen.yaml api/buf.{service-name}.openapi.gen.yaml
# Edit the config to point to your service's proto files
```

1. Create proto files:
   - HTTP interfaces: `api/protos/{service-name}/service/v1/i_*.proto`
   - gRPC interfaces: `api/protos/{domain}/service/v1/{domain}.proto`

### Proto File Organization

Following go-wind-admin conventions:

- **i_*.proto files**: HTTP interfaces with `google.api.http` annotations
  - Located in: `api/protos/{service-name}/service/v1/`
  - Package: `{service-name}.service.v1`
  - Used for: OpenAPI generation

- **Other .proto files**: Pure gRPC interfaces
  - Located in: `api/protos/{domain}/service/v1/`
  - Package: `{domain}.service.v1`
  - Used for: gRPC service implementation

Example:

```
api/protos/
├── krathub/service/v1/
│   ├── i_auth.proto      # HTTP: package krathub.service.v1
│   ├── i_user.proto      # HTTP: package krathub.service.v1
│   └── i_test.proto      # HTTP: package krathub.service.v1
├── auth/service/v1/
│   └── auth.proto        # gRPC: package auth.service.v1
├── user/service/v1/
│   └── user.proto        # gRPC: package user.service.v1
└── test/service/v1/
    └── test.proto        # gRPC: package test.service.v1
```

## Makefile Targets

### Root Makefile

| Target | Description |
|--------|-------------|
| `make init` | Initialize development environment |
| `make plugin` | Install protoc plugins |
| `make cli` | Install CLI tools |
| `make api` | Generate all protobuf Go code |
| `make openapi` | Generate OpenAPI docs for all services |
| `make wire` | Generate wire code for all services |
| `make gen` | Generate all code (wire + api + openapi) |
| `make build` | Build all services |
| `make build_only` | Build without code generation |
| `make all` | Generate and build all services |
| `make test` | Run tests |
| `make lint` | Run linter |
| `make clean` | Clean build artifacts |

### Service Makefile (app/{service}/service/)

All services use the common `app.mk` file:

| Target | Description |
|--------|-------------|
| `make build` | Build service (with code generation) |
| `make build_only` | Build service only |
| `make run` | Run service |
| `make api` | Generate protobuf code |
| `make openapi` | Generate OpenAPI documentation |
| `make wire` | Generate wire code |
| `make gen` | Generate all code |
| `make clean` | Clean build files |
| `make env` | Show environment variables |

## Code Generation

### Protobuf Generation

```bash
# Generate all protobuf code
make api

# Or use buf directly
cd api && buf generate
```

### OpenAPI Generation

Each service has its own OpenAPI config file:

- `api/buf.krathub.openapi.gen.yaml` - Krathub service

OpenAPI files are generated to each service directory:

- `app/krathub/service/openapi.yaml`

```bash
# Generate OpenAPI for all services
make openapi

# Generate for specific service
cd app/krathub/service && make openapi
```

### Wire Generation

```bash
# Generate wire code for all services
make wire

# Generate for specific service
cd app/krathub/service && make wire
```

## Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

Variables:

- `PROJECT_NAME`: Project name (default: krathub)
- `SERVICE_APP_VERSION`: Service version (default: 0.0.1)
- `VERSION`: Build version (overridden by git tags)

## Tools Required

- Go 1.21+
- buf CLI
- protoc plugins:
  - protoc-gen-go
  - protoc-gen-go-grpc
  - protoc-gen-go-http
  - protoc-gen-go-errors
  - protoc-gen-openapi
  - protoc-gen-validate
- wire
- golangci-lint

Install all tools:

```bash
make init
```

## Tips

1. **Always run from root directory**: The Makefile is designed to work from the project root
2. **Use `make gen` before building**: Ensures all generated code is up-to-date
3. **Check `make env`**: Verify environment variables in each service
4. **OpenAPI location**: Each service's OpenAPI file is in its own directory

## References

- [go-wind-admin](https://github.com/tx7do/go-wind-admin) - Project structure reference
- [Buf](https://buf.build/) - Protocol buffer tooling
- [Kratos](https://go-kratos.dev/) - Microservice framework
- [Wire](https://github.com/google/wire) - Dependency injection
