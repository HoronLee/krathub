# Krathub

[简体中文](README.md) | English

A production-ready, feature-rich template for building microservices with the [Kratos](https://github.com/go-kratos/kratos) framework.

Krathub aims to provide a best-practice project layout that includes many commonly used components, allowing developers to focus on business logic from day one.

## Features

*   **Monorepo Layout:** Clear and scalable project structure.
*   **API-First:** Define APIs with Protobuf and generate code for gRPC, HTTP, and validation.
*   **Database:** Integrated with GORM and provides `make gendb` for easy ORM code generation.
*   **Dependency Injection:** Using `google/wire` for compile-time dependency injection.
*   **Configuration:** Centralized configuration management with support for Consul and Nacos.
*   **Service Discovery:** Service registration and discovery with Consul and Nacos.
*   **Authentication:** JWT-based authentication middleware.
*   **Containerized:** Comes with `Dockerfile` and `docker-compose.yml` for easy deployment.
*   **Tooling:** A comprehensive `Makefile` for development and CI/CD.

## Getting Started

### Prerequisites

*   [Go](https://golang.org/doc/install) >= 1.21
*   [Protocol Buffers](https://grpc.io/docs/protoc-installation/)
*   [Kratos Tool](https://go-kratos.dev/docs/getting-started/start#install-kratos)
*   Make

### Installation & Usage

1.  **Create a new project from this template:**

    ```bash
    kratos new your-project-name -r https://github.com/HoronLee/krathub.git
    ```

2.  **Enter the project directory:**

    ```bash
    cd your-project-name
    ```

3.  **Initialize the development environment:**

    This command will install all the necessary Go tools for code generation.

    ```bash
    make init
    ```

4.  **Configure your application:**

    Modify `configs/config.yaml` to set up your database, Redis, and other services.

5.  **Generate database code:**

    After configuring your database connection, you can generate GORM models from your database schema.

    ```bash
    make gendb
    ```

6.  **Generate API and other code:**

    ```bash
    make proto
    make wire
    ```

7.  **Run the application:**

    ```bash
    make run
    ```

## Makefile Commands

The `Makefile` provides a set of useful commands to streamline development:

| Command      | Description                                               |
|--------------|-----------------------------------------------------------|
| `make init`    | Install required Go tools for code generation.            |
| `make proto`   | Generate all proto files (API, config, validate, errors). |
| `make gendb`   | Generate GORM models from the database.                   |
| `make wire`    | Run `wire` to generate dependency injection code.         |
| `make build`   | Build the application binary.                             |
| `make run`     | Run the application.                                      |
| `make test`    | Run tests.                                                |
| `make fmt`     | Format the source code.                                   |
| `make clean`   | Remove all generated files.                               |
| `make all`     | A shortcut to run `clean`, `proto`, `gendb`, `wire`, `generate`, and `build`. |
| `make help`    | Show all available commands.                              |

## Project Structure

The project follows the standard Kratos layout, with a few additions:

```
.
├── api         # Protobuf API definitions
├── bin         # Compiled binaries
├── cmd         # Main application entrypoints
├── configs     # Configuration files
├── internal    # Core application logic
│   ├── biz     # Business logic (usecases)
│   ├── client  # gRPC/HTTP clients for external services
│   ├── conf    # Protobuf definitions for configuration
│   ├── consts  # Constants
│   ├── data    # Data access layer (repositories)
│   ├── server  # gRPC and HTTP servers
│   └── service # Service layer (connects transport to business logic)
├── manifest    # Deployment manifests (Docker, K8s)
├── pkg         # Shared libraries
└── third_party # Third-party Protobuf files
```

## Configuration

Configuration is managed via `configs/config.yaml`. The application uses a hierarchical configuration system that can be extended with configuration centers like Consul or Nacos. See `internal/conf` for the Protobuf definitions of the configuration structure.

## API Development Flow

1.  **Define:** Create or modify `.proto` files in the `api/` directory.
2.  **Generate:** Run `make proto` to generate gRPC, HTTP, validation, and error code.
3.  **Implement Service:** Implement the service interface in `internal/service/`.
4.  **Implement Business Logic:** Implement the core logic in `internal/biz/`.
5.  **Implement Data Access:** Implement the data access logic in `internal/data/`.
6.  **Wire it up:**
    *   Add dependencies to the `ProviderSet` in the appropriate `wire.go` file (`data/data.go`, `biz/biz.go`, etc.).
    *   Run `make wire` to generate the dependency injection code.
7.  **Run:** Run `make run` to start the server and test your new endpoint.
