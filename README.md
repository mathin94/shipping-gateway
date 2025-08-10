# Shipping Gateway

Shipping Gateway project for integration with third party shipment services RESTful APIs in Go using [Gin](https://github.com/gin-gonic/gin), following Clean Architecture principles.

---

## Project Structure

```
shipping-gateway/
├── cmd/
│   ├── web/                    # Entry point for the web server
│   └── worker/                 # Entry point for background workers or jobs
├── config.yaml                 # Application configuration file
├── db/
│   └── migration/              # Database migration scripts
├── internal/
│   ├── config/                 # Configuration setup for frameworks and libraries
│   │   ├── app.go
│   │   ├── gin.go
│   │   ├── gorm.go
│   │   ├── logrus.go
│   │   ├── validator.go
│   │   └── viper.go
│   ├── delivery/
│   │   ├── http/
│   │   │   ├── health_check_controller.go  # Health check endpoint
│   │   │   ├── middleware/                # HTTP middleware
│   │   │   └── route/
│   │   │       └── route.go               # Route definitions
│   │   └── messaging/                     # Messaging delivery (e.g., consumers, producers)
│   ├── entity/                            # Domain entities (business models)
│   ├── gateway/
│   │   └── messaging/                     # Messaging gateway implementations
│   ├── model/
│   │   └── converter/                     # Data model converters
│   ├── repository/
│   │   └── repository.go                  # Data access interfaces and implementations
│   └── usecase/                           # Business logic layer (use cases/interactors)
├── test/                                  # Test files and test data
├── go.mod                                 # Go module definition
├── go.sum                                 # Go module checksums
├── LICENSE                                # License file
└── README.md                              # Project documentation
```

---

## Directory Overview

| Directory / File                  | Description                                                                                       |
|-----------------------------------|---------------------------------------------------------------------------------------------------|
| `cmd/web/`                        | Entry point for the web server application.                                                       |
| `cmd/worker/`                     | Entry point for background workers or scheduled jobs.                                             |
| `config.yaml`                     | Centralized application configuration file.                                                       |
| `db/migration/`                   | Database migration scripts for schema management.                                                 |
| `internal/config/`                | Configuration setup for frameworks and libraries (Gin, Gorm, Logrus, Viper, Validator).           |
| `internal/delivery/http/`         | HTTP delivery layer: controllers, middleware, and route definitions.                              |
| `internal/delivery/messaging/`    | Messaging delivery (e.g., consumers and producers for message queues).                            |
| `internal/entity/`                | Domain entities representing core business objects.                                               |
| `internal/gateway/messaging/`     | Messaging gateway implementations (e.g., for brokers/integrations).                               |
| `internal/model/converter/`       | Data model converters for mapping between layers.                                                 |
| `internal/repository/`            | Data access interfaces and implementations.                                                       |
| `internal/usecase/`               | Business logic and application use cases.                                                         |
| `test/`                           | Test code and related resources.                                                                  |
| `go.mod` / `go.sum`               | Go module files for dependency management.                                                        |
| `LICENSE`                         | Project license.                                                                                  |
| `README.md`                       | Project documentation and instructions.                                                           |

---

## Clean Architecture Principles

This project follows **Clean Architecture** to promote separation of concerns and facilitate maintainability:

- **Delivery Layer**: Handles incoming requests (HTTP controllers, middleware, routing) and messaging consumers.
- **Usecase Layer**: Contains business logic and application-specific rules.
- **Repository Layer**: Abstracts data access for infrastructure independence.
- **Entity Layer**: Defines core business models/entities.
- **Gateway Layer**: Handles external integrations (e.g., messaging systems).

---

## Getting Started

1. **Clone the repository**
   ```sh
   git clone https://github.com/mathin94/shipping-gateway.git
   cd shipping-gateway
   ```

2. **Install dependencies**
   ```sh
   go mod tidy
   ```

3. **Configure the application**
    - Edit `config.yaml` to match your environment.

4. **Run database migrations**
   ```sh
   # Use your migration tool of choice (e.g., migrate, goose, etc.)
   ```

5. **Start the web server**
   ```sh
   go run cmd/web/main.go
   ```

6. **(Optional) Start the worker**
   ```sh
   go run cmd/worker/main.go
   ```

---

## License

This project is licensed under the [MIT License](LICENSE).

---

> **Inspired by [Golang Clean Architecture](https://github.com/khannedy/golang-clean-architecture)**