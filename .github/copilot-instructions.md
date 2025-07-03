# GitHub Copilot Custom Instructions for Go Backend Service

## Role Instructions
You are a Go backend service developer in globalxtreme. You are responsible for implementing boilerplate and helping developers follow the coding standards and architecture patterns established in the project. You will provide code suggestions, explanations, and best practices for building scalable and maintainable backend services. Everytime you suggest code, you must follow the existing architecture patterns, use the established error handling mechanisms, leverage the custom framework utilities, maintain consistency with existing code style, consider performance and scalability, include proper validation and error handling, and use appropriate logging where necessary. Ensure always use core library for common operations and utilities which can be found in `github.com/globalxtreme/go-core/v2`.

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
      - **grpc/**: gRPC saga implementations
      - **privateapi/**: Private API saga implementations
    - **thirdparty/**: Integrations with third-party api services such as Email, Telegram, etc.
    - ***/**: Domain-specific packages, hooks and helpers
  - ***/**: Domain-specific code that only applies to the specific domain
    - **repository/**: Data access layer for the domain
    - **service/**: Business logic for the domain
    - **mail/**: Email sending and templating
    - **excel/**: Excel file processing and generation
    - **pdf/**: PDF generation
- **storage/**: File storage and logs
  - **app/**: Application-specific storage
  - **logs/**: Application logs
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
import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

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

### Database Migrations
- When creating new migrations, *YOU MUST ALWAYS* use cli command: `go run ./cmd/main.go gen:migration <NAME>`
- For migration that creates new tables, use `<FEATURE_NAME>` as the `<NAME>`, e.g., `Product`
- For migration that modifies existing tables, use `<FEATURE_NAME>Batch<NumberOfTimesTheFeatureWasModified>` format. Batch number for migration that modifies existing tables should start from 2, make sure to look for the last batch number used in the file in `internal/app/database/migration/`
- Then register migrations in `internal/app/database/migrate.go`
- After new migration is created, you should implement the file content following the example below.

Example migration for new table:
```go
package migration

import (
	"os"
	"service/internal/pkg/config"
	"service/internal/pkg/model"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type Product_1726651240922259 struct{}

func (Product_1726651240922259) Reference() string {
	return "Product_1726651240922259"
}

func (Product_1726651240922259) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")

	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: model.Product{}, Owner: owner},
		{Connection: config.PgSQL, CreateTable: model.ProductComponentCategory{}, Owner: owner},
		{Connection: config.PgSQL, CreateTable: model.ProductVariantInformation{}, Owner: owner},
	}
}

func (Product_1726651240922259) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}

```

Example migration for modifying existing table:
```go
package migration

import (
	"os"
	"service/internal/pkg/config"
	"service/internal/pkg/model"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type ProductBatch2_1738380774908382 struct{}

func (ProductBatch2_1738380774908382) Reference() string {
	return "ProductBatch2_1738380774908382"
}

func (ProductBatch2_1738380774908382) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")

	return []xtremedb.Table{
		{Connection: config.PgSQL, RenameTable: xtremedb.Rename{Old: "setting_add_ons", New: "product_component_add_ons"}, Owner: owner},
		{Connection: config.PgSQL, RenameTable: xtremedb.Rename{Old: "setting_informations", New: "product_component_informations"}, Owner: owner},
    {Connection: config.PgSQL, DropTable: model.ProductVariant{}, Owner: owner},
	}
}

func (ProductBatch2_1738380774908382) Columns() []xtremedb.Column {
	return []xtremedb.Column{
		{
			Connection:  config.PgSQL,
			Model:       model.ProductVariant{},
			DropColumns: []string{"column1", "column2"},
			RenameColumns: []xtremedb.Rename{
				{
					Old: "oldColumnName",
					New: "newColumnName",
				},
			},
			AddColumns: []string{"newColumn1", "newColumn2"},
			AlterColumns: []string{"column3", "column4"},
		},
	}
}

```

### Database Seeder
- Put in `internal/app/database/seeder/`
- Implement seeders for initial data population
- Register seeders in `internal/app/database/seeder.go`

Example seeder:
```go
package seeder

import (
	"service/internal/pkg/config"
	"service/internal/pkg/model"
)

type TestingSeeder struct{}

func (seed *TestingSeeder) Seed() {
	testings := seed.setTestingData()
	for _, testing := range testings {
		var count int64
		config.PgSQL.Model(&model.Testing{}).Where("name = ?", testing["name"]).Count(&count)
		if count > 0 {
			continue
		}

		config.PgSQL.Create(&model.Testing{
			Name: testing["name"].(string),
		})
	}
}

func (seed *TestingSeeder) setTestingData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name": "Testing",
		},
	}
}
```

### Form Request
- Use `internal/pkg/form/` for request validation
- Implement form validation using `github.com/go-playground/validator/v10`

Example form structure:
```go
package form

import (
	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
	"net/http"
	"service/internal/pkg/core"
)

type TestingForm struct { // Example form structure
	Name string   `json:"name"`
	Subs []string `json:"subs" validate:"required"`
}

func (rule *TestingForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(rule)
}

func (rule *TestingForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &rule) //if content type is application/json
  rule.Request = r // if content type is multipart/form-data
}
```

### Response Parser
- Use `internal/pkg/parser/` for data transformation
- Implement parsers for API responses

Example parser structure:
```go
package parser

import (
	xtremefs "github.com/globalxtreme/go-core/v2/filesystem"
	"service/internal/pkg/model"
)

type TestingParser struct {
	Array  []model.Testing
	Object model.Testing
}

func (parser TestingParser) Get() []interface{} { //Get all data
	var result []interface{}

	for _, activity := range parser.Array {
		firstParser := TestingParser{Object: activity}
		result = append(result, firstParser.First())
	}

	return result
}

func (parser TestingParser) First() interface{} { // Get single data
	activity := parser.Object

	var resSubs []interface{}
	for _, sub := range activity.Subs {
		resSubs = append(resSubs, map[string]interface{}{
			"id":        sub.ID,
			"name":      sub.Name,
			"createdAt": sub.CreatedAt.Format("02/01/2006 15:04"),
		})
	}

	return map[string]interface{}{
		"id":        activity.ID,
		"name":      activity.Name,
		"createdAt": activity.CreatedAt.Format("02/01/2006 15:04"),
		"file":      xtremefs.Storage{}.GetFullPathURL("ckH2cahaAaDMNVgS2xdM1697957810885349000.png"),
		"subs":      resSubs,
	}
}
```

### Domain Specific Codes
- Place domain-specific code in `internal/<domain>/`
- Use `repository/` for data access layer

Example repository structure:
```go
package repository

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	"gorm.io/gorm"
	"net/url"
	"service/internal/pkg/config"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form"
	"service/internal/pkg/model"
)

type TestingRepository interface {
	core.TransactionRepository

	FirstById(id any, args ...func(query *gorm.DB) *gorm.DB) model.Testing
	Find(parameter url.Values) ([]model.Testing, interface{}, error)

	Store(form form.TestingForm) model.Testing
	Delete(testing model.Testing)

	AddSub(testing model.Testing, sub string) model.TestingSub
	DeleteSub(testingSub model.TestingSub)
}

func NewTestingRepository(args ...*gorm.DB) TestingRepository {
	repository := testingRepository{}
	if len(args) > 0 {
		repository.transaction = args[0]
	}

	return &repository
}

type testingRepository struct {
	transaction *gorm.DB
}

func (repo *testingRepository) SetTransaction(tx *gorm.DB) {
	repo.transaction = tx
}

func (repo *testingRepository) FirstById(id any, args ...func(query *gorm.DB) *gorm.DB) model.Testing {
	var testing model.Testing

	query := config.PgSQL
	if len(args) > 0 {
		query = args[0](query)
	}

	err := query.First(&testing, "id = ?", id).Error
	if err != nil {
		error2.ErrXtremeTestingGet(err.Error())
	}

	return testing
}

func (repo *testingRepository) Find(parameter url.Values) ([]model.Testing, interface{}, error) {
	fromDate, toDate := core.SetDateRange(parameter)

	query := config.PgSQL.Preload("Subs").
		Where("\"createdAt\" BETWEEN ? AND ?", fromDate, toDate)

	if search := parameter.Get("search"); len(search) > 3 {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	testings, pagination, err := xtrememodel.Paginate(query.Order("id DESC"), parameter, model.Testing{})
	if err != nil {
		return nil, nil, err
	}

	return testings, pagination, nil
}

func (repo *testingRepository) Store(form form.TestingForm) model.Testing {
	testing := model.Testing{
		Name: form.Name,
	}

	err := repo.transaction.Create(&testing).Error
	if err != nil {
		error2.ErrXtremeTestingSave(err.Error())
	}

	return testing
}

func (repo *testingRepository) Delete(testing model.Testing) {
	err := repo.transaction.Delete(&testing).Error
	if err != nil {
		error2.ErrXtremeTestingDelete(err.Error())
	}
}

func (repo *testingRepository) AddSub(testing model.Testing, sub string) model.TestingSub {
	testingSub := model.TestingSub{
		TestingId: testing.ID,
		Name:      sub,
	}

	err := repo.transaction.Create(&testingSub).Error
	if err != nil {
		error2.ErrXtremeTestingSubSave(err.Error())
	}

	return testingSub
}

func (repo *testingRepository) DeleteSub(testingSub model.TestingSub) {
	err := repo.transaction.Delete(&testingSub).Error
	if err != nil {
		error2.ErrXtremeTestingSubDelete(err.Error())
	}
}

```

- Use `service/` for business logic
  - Every service that would modify the data (Create/Update/Delete) should use a transaction and save their action as `Activity`
Example service structure:
```go
package service

import (
	"fmt"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	error2 "service/internal/pkg/error"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/model"
	"service/internal/pkg/port"
	"service/internal/testing/repository"

	xtremefs "github.com/globalxtreme/go-core/v2/filesystem"
	"gorm.io/gorm"
)

type TestingService interface {
	SetTransaction(tx *gorm.DB)
	SetActivityRepository(repo port.ActivityRepository)

	Create(form form2.TestingForm) model.Testing
	Update(form form2.TestingForm, id int64) model.Testing
	Delete(id int64) error
	UploadByFile(form form2.TestingUploadForm) map[string]interface{}
	UploadByContent(form form2.TestingUploadContentForm) map[string]interface{}
}

func NewTestingService() TestingService {
	return &testingService{}
}

type testingService struct {
	tx *gorm.DB

	repository         repository.TestingRepository
	activityRepository port.ActivityRepository
}

func (srv *testingService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *testingService) SetActivityRepository(repo port.ActivityRepository) {
	srv.activityRepository = repo
}

func (srv *testingService) Create(form form2.TestingForm) model.Testing {
	var testing model.Testing

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewTestingRepository(tx)

		testing = srv.repository.Store(form)

		for _, sub := range form.Subs {
			testingSub := srv.repository.AddSub(testing, sub)
			testing.Subs = append(testing.Subs, testingSub)
		}

		activity.UseActivity{}.SetReference(testing).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Enter new testing: %s [%d]", testing.Name, testing.ID))

		return nil
	})

	return testing
}

func (srv *testingService) Update(form form2.TestingForm, id int64) model.Testing {
	var testing model.Testing

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewTestingRepository(tx)

		testing = srv.repository.FirstById(id, func(query *gorm.DB) *gorm.DB {
			return query.Preload("Subs")
		})
		if testing.ID == 0 {
			error2.ErrXtremeTestingGet("Testing not found")
		}

		testing.Name = form.Name

		for _, sub := range testing.Subs {
			srv.repository.DeleteSub(sub)
		}

		for _, sub := range form.Subs {
			testingSub := srv.repository.AddSub(testing, sub)
			testing.Subs = append(testing.Subs, testingSub)
		}

		activity.UseActivity{}.SetReference(testing).SetNewProperty(constant.ACTION_UPDATE).
			Save(fmt.Sprintf("Update testing: %s [%d]", testing.Name, testing.ID))

		return nil
	})

	return testing
}

func (srv *testingService) Delete(id int64) error {
	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewTestingRepository(tx)
		testing := srv.repository.FirstById(id, func(query *gorm.DB) *gorm.DB {
			return query.Preload("Subs")
		})
		if testing.ID == 0 {
			error2.ErrXtremeTestingGet("Testing not found")
		}
		srv.repository.Delete(testing)
		activity.UseActivity{}.SetReference(testing).SetNewProperty(constant.ACTION_DELETE).
			Save(fmt.Sprintf("Delete testing: %s [%d]", testing.Name, testing.ID))
		return nil
	})
	return nil
}

func (srv *testingService) UploadByFile(form form2.TestingUploadForm) map[string]interface{} {
	uploader := xtremefs.Uploader{Path: constant.PathImageTesting(), IsPublic: true}
	filePath, err := uploader.MoveFile(form.Request, "testFile[testing][0]")
	if err != nil {
		error2.ErrXtremeTestingSave("Unable to upload file: " + err.Error())
	}

	storage := xtremefs.Storage{IsPublic: uploader.IsPublic}

	return map[string]interface{}{
		"url":      storage.GetFullPathURL(filePath.(string)),
		"fullPath": storage.GetFullPath(filePath.(string)),
		"path":     filePath.(string),
	}
}

func (srv *testingService) UploadByContent(form form2.TestingUploadContentForm) map[string]interface{} {
	uploader := xtremefs.Uploader{Path: constant.PathImageTesting(), IsPublic: true}
	filePath, err := uploader.MoveContent(form.Content)
	if err != nil {
		error2.ErrXtremeTestingSave("Unable to upload file: " + err.Error())
	}

	storage := xtremefs.Storage{IsPublic: uploader.IsPublic}

	return map[string]interface{}{
		"url":      storage.GetFullPathURL(filePath.(string)),
		"fullPath": storage.GetFullPath(filePath.(string)),
		"path":     filePath.(string),
	}
}
```

- Use `mail/` for email sending and templating
- Use `excel/` for Excel file processing and generation
- Use `pdf/` for PDF generation



### Route Handlers
- When creating new handlers, *YOU MUST ALWAYS* use cli command: `go run ./cmd/main.go gen:handler <NAME> --type=<web/mobile> --resource`
- Use `--type=web` for web API handlers, `--type=mobile` for mobile API handlers
- Use `--resource` flag to generate resourceful handlers
- Implement handlers in `internal/app/api/<type>/handler/`
- Register created handlers in `internal/app/api/<type>/router.go`

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
