# GitHub Copilot Custom Instructions for Go Backend Service

## Project Overview
This is a Go backend service built with a layered architecture following Domain-Driven Design principles. The service uses gRPC, REST APIs, message queues (RabbitMQ), and various third-party integrations.

## Architecture & Patterns

### Project Structure
- **cmd/**: Entry points and CLI commands using Cobra
- **internal/**: Private application code
  - **app/**: Application layer (API routes, gRPC servers, console commands, queues). Entry points for the application.
    - **api/**: HTTP API handlers
      - **router.go**: Main router for API endpoints
      - ***/** : Specific API types (web, mobile, private)
        - **router.go**: Router for specific API type
        - **handler/**: Handlers for API endpoints
    - **console/**: CLI commands
      - **command/**: Command implementations
      - **kernel.go**: Command registration
    - **database/**: Database migrations and seeding
      - **migration/**: Database migration files
      - **seeder/**: Database seeder files
      - **migrate.go**: Migration registration
      - **seeder.go**: Seeder registration
    - **grpc/**: gRPC server implementations
      - **server/**: gRPC logic handlers
      - **worker.go**: gRPC worker registration
    - **privateapi/**: Private API handlers
      - **handler/**: Handlers for private API endpoints
      - **router.go**: Router for private API endpoints
    - **queue/**: Internal message queue handlers using redis
      - **job/**: Queue job implementations
      - **worker.go**: Queue worker registration
    - **rabbitmq/**: RabbitMQ message queue handlers
  - **pkg/**: Domain/business logic packages
    - **config/**: Configuration management
    - **constant/**: Application constants
    - **core/**: Core utilities and helpers
    - **encrypt/**: Encryption and decryption utilities
    - **error/**: Custom error types and handling
    - **form/**: Interfaces for form transmission
    - **grpc/**: gRPC generated code and utilities
    - **layout/**: Layout templates for HTML rendering for emails and pdfs
    - **middleware/**: HTTP middleware for accessing services
    - **model/**: Database models using GORM
    - **number/**: Number utilities
    - **parser/**: Data transformation utilities usually for API responses
    - **port/**: Repository interfaces and ports
    - **saga/**: Saga pattern implementations for distributed transactions
    - **service/**: Business logic services
    - ***/**: Domain-specific utilities, hooks and helpers
  - ***/**: Domain-specific packages
  - 
- **storage/**: File storage and logs
- **stubs/**: Code generation templates

### Layered Architecture
Follow this pattern when creating new features:
1. **Model** (`internal/pkg/model/`): Database entities using GORM
2. **Repository** (`internal/*/repository/`): Data access layer
3. **Service** (`internal/*/service/`): Business logic
4. **Handler** (`internal/app/*/handler/`): HTTP/gRPC request handlers
5. **Parser** (`internal/pkg/parser/`): Data transformation
6. **Form** (`internal/pkg/form/`): Request validation

## Coding Standards

### Go Conventions
- Use PascalCase for exported functions, types, and methods
- Use camelCase for unexported functions and variables
- Follow standard Go project layout
- Use meaningful package names that reflect their purpose

### Database Models
- Extend `xtrememodel.BaseModel` for all entities
- Use GORM tags for database mapping
- Implement `TableName()` method for custom table names
- Implement `SetReference()` method for reference handling

Example model structure:
```go
type YourModel struct {
    xtrememodel.BaseModel
    Name string `gorm:"column:name;type:varchar(250);default:null"`
}

func (YourModel) TableName() string {
    return "your_table"
}

func (model YourModel) SetReference() uint {
    return model.BaseModel.ID
}
```

### Repository Pattern
- Create interfaces in `internal/pkg/port/repository.go`
- Implement repositories in respective domain folders
- Use GORM for database operations
- Handle errors appropriately

### Service Layer
- Business logic should be in service layer
- Services should depend on repository interfaces
- Handle validation and business rules
- Use appropriate error handling from `internal/pkg/error/`

### API Handlers
- Separate handlers for different API types (mobile, web, private)
- Use form validation for request validation
- Return consistent response formats
- Handle errors using the error package

## Technology Stack

### Core Dependencies
- **GORM**: ORM for database operations
- **Gorilla Mux**: HTTP router
- **gRPC**: For service communication
- **Protocol Buffers**: Message serialization
- **RabbitMQ**: Message queuing with amqp091-go
- **Cobra**: CLI framework
- **Go-cache**: In-memory caching
- **Excelize**: Excel file processing
- **CORS**: Cross-origin resource sharing

### Custom Framework
- Uses `github.com/globalxtreme/go-core/v2` as base framework
- Leverage framework utilities for common operations

## File Generation Patterns

### When creating new features:
1. **Model**: Create in `internal/pkg/model/`
2. **Migration**: Use timestamp-based naming in `internal/app/database/migration/`
3. **Repository**: Interface in `port/repository.go`, implementation in domain folder
4. **Service**: Business logic in `internal/testing/service/` (adjust domain as needed)
5. **Forms**: Request validation in `internal/pkg/form/`
6. **Handlers**: API endpoints in `internal/app/api/*/handler/`
7. **Parsers**: Data transformation in `internal/pkg/parser/`

### File Naming Conventions
- Use PascalCase for Go files: `YourModelName.go`
- Include purpose in filename: `YourModelRepository.go`, `YourModelService.go`
- Use descriptive names for handlers: `YourFeatureHandler.go`

## Common Patterns

### Error Handling
- Use custom error types from `internal/pkg/error/`
- Create domain-specific errors
- Handle and log errors appropriately

### Configuration
- Configuration files in `internal/pkg/config/`
- Use environment-based configuration
- Follow existing patterns for new config types

### Constants
- Define constants in `internal/pkg/constant/`
- Group related constants in separate files
- Use descriptive constant names

### Queue Jobs
- Implement jobs in `internal/app/queue/job/`
- Follow existing job patterns
- Handle job failures gracefully

### gRPC Services
- Protocol buffers in `cmd/protobuf/`
- Generated code in `internal/pkg/grpc/`
- Server implementations in `internal/app/grpc/server/`

## Testing Patterns
- Create test files with `_test.go` suffix
- Use table-driven tests where appropriate
- Mock external dependencies
- Test business logic thoroughly

## CLI Commands
- Use Cobra for CLI commands
- Implement commands in `internal/app/console/command/`
- Follow existing command patterns

## Best Practices
1. Always validate input data using forms
2. Use proper error handling and logging
3. Follow the established repository pattern
4. Keep business logic in service layer
5. Use appropriate HTTP status codes
6. Implement proper transaction handling for database operations
7. Use constants for magic numbers and strings
8. Follow Go naming conventions
9. Write comprehensive tests
10. Document complex business logic

## When suggesting code:
- Follow the existing architecture patterns
- Use the established error handling mechanisms
- Leverage the custom framework utilities
- Maintain consistency with existing code style
- Consider performance and scalability
- Include proper validation and error handling
- Use appropriate logging where necessary
