# GitHub Copilot Custom Instructions for Go Backend Service

## Role Instructions
You are a Go backend service developer in globalxtreme. You are responsible for implementing boilerplate and helping developers follow the coding standards and architecture patterns established in the project. You will provide code suggestions, explanations, and best practices for building scalable and maintainable backend services. Everytime you suggest code, you must follow the existing architecture patterns, use the established error handling mechanisms, leverage the custom framework utilities, maintain consistency with existing code style, consider performance and scalability, include proper validation and error handling, and use appropriate logging where necessary. Ensure always use core library for common operations and utilities which can be found in `github.com/globalxtreme/go-core/v2`.

***"Important"*** Don't leave any comments in the code you suggest, unless it is a comment that is already in the codebase. Do not add any comments that are not already in the codebase.

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
	"service/internal/pkg/model"

	"gorm.io/gorm"
)

type CakeIngredientSeeder struct {}

func (seed *CakeIngredientSeeder) Seed() {
	ingredients := seed.setIngredientsData()
	for _, ingredient := range ingredients {
		var count int64
		config.PgSQL.Model(&model.CakeComponentIngredient{}).Where("name = ?", ingredient["name"]).Count(&count)
		if count > 0 {
			continue
		}

		config.PgSQL.Model(&model.CakeComponentIngredient{}).Where("name = ?", ingredient["name"]).Count(&count)Create(&model.CakeComponentIngredient{
			Name:        ingredient["name"].(string),
			Description: ingredient["description"].(string),
			UnitPrice:   ingredient["unit_price"].(float64),
			Unit:        ingredient["unit"].(string),
		})
	}
}

// --- UNEXPORTED FUNCTIONS ---

func (seed *CakeIngredientSeeder) setIngredientsData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "Eggs",
			"description": "Fresh chicken eggs",
			"unit_price":  2000.0,
			"unit":        "pcs",
		},
		{
			"name":        "Flour",
			"description": "All-purpose wheat flour",
			"unit_price":  10000.0,
			"unit":        "kg",
		},
		{
			"name":        "Sugar",
			"description": "Granulated white sugar",
			"unit_price":  12000.0,
			"unit":        "kg",
		},
		{
			"name":        "Chocolate",
			"description": "Chocolate compound for baking",
			"unit_price":  25000.0,
			"unit":        "kg",
		},
	}
}

```

### Form Request
- Use `internal/pkg/form/` for request validation
- Implement form validation using `github.com/go-playground/validator/v10`
- Json attributes should be in camelCase, e.g., `name` instead of `Name`, `sellPrice` instead of `sell_price`

Example form structure:
```go
package form

import (
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type CakeForm struct {
	Name        string                   `json:"name" validate:"required,max=250"`
	Description string                   `json:"description"`
	Margin      float64                  `json:"margin" validate:"required,gte=0"`
	Unit        string                   `json:"unit" validate:"max=50"`
	Stock       int                      `json:"stock" validate:"gte=0"`
	Ingredients []CakeCompIngredientForm `json:"ingredients" validate:"required,dive"`
	Costs       []CakeCompCostForm       `json:"costs" validate:"required,dive"`
}

type CakeCompIngredientForm struct {
	IngredientID uint    `json:"ingredientId" validate:"required,gt=0"`
	Amount       float64 `json:"amount" validate:"required,gte=0"`
	Unit         string  `json:"unit" validate:"required,max=50"`
}

type CakeCompCostForm struct {
	CostType string  `json:"type" validate:"required,max=100"`
	Cost     float64 `json:"cost" validate:"required,gte=0"`
}

func (f *CakeForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(f)
}

func (f *CakeForm) APIParse(r *http.Request) {
	form := core.BaseForm{}
	form.APIParse(r, &f)
}

```

### Response Parser
- Use `internal/pkg/parser/` for data transformation
- Implement parsers for API responses
- Returned attributes should be in camelCase `sellPrice` instead of `sell_price`

Example parser structure:
```go
package parser

import (
	"service/internal/pkg/model"
)

type CakeParser struct {
	Array  []model.Cake
	Object model.Cake
}

func (parser CakeParser) CreateActivity(action string) interface{} {
	return parser.Brief()
}

func (parser CakeParser) DeleteActivity(action string) interface{} {
	return parser.Brief()
}

func (parser CakeParser) GeneralActivity(action string) interface{} {
	return parser.Brief()
}

func (parser CakeParser) UpdateActivity(action string) interface{} {
	return parser.Brief()
}

func (parser CakeParser) Briefs() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, CakeParser{Object: obj}.Brief())
	}
	return result
}

func (parser CakeParser) Brief() interface{} {
	cakeObj := parser.Object
	var recipes []interface{}
	for _, recipe := range cakeObj.Recipes {
		recipes = append(recipes, CakeRecipeParser{Object: recipe}.First())
	}
	var costs []interface{}
	for _, cost := range cakeObj.Costs {
		costs = append(costs, CakeCostParser{Object: cost}.First())
	}
	return map[string]interface{}{
		"id":          cakeObj.ID,
		"name":        cakeObj.Name,
		"description": cakeObj.Description,
		"margin":      cakeObj.Margin,
		"sellPrice":   cakeObj.SellPrice,
		"unit":        cakeObj.Unit,
		"stock":       cakeObj.Stock,
		"createdAt":   cakeObj.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":   cakeObj.UpdatedAt.Format("02/01/2006 15:04"),
		"recipes":     recipes,
		"costs":       costs,
	}
}

func (parser CakeParser) Get() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, CakeParser{Object: obj}.First())
	}
	return result
}

func (parser CakeParser) First() interface{} {
	cakeObj := parser.Object
	return map[string]interface{}{
		"id":          cakeObj.ID,
		"name":        cakeObj.Name,
		"description": cakeObj.Description,
		"margin":      cakeObj.Margin,
		"sellPrice":   cakeObj.SellPrice,
		"unit":        cakeObj.Unit,
		"stock":       cakeObj.Stock,
		"createdAt":   cakeObj.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":   cakeObj.UpdatedAt.Format("02/01/2006 15:04"),
	}
}

```

### Domain Specific Codes
- Place domain-specific code in `internal/<domain>/`
- Domain cannot interact directly with other domains, use either `internal/pkg/` or expose demanded repository/service function interfaces inside `internal/pkg/port`

- Use `repository/` for data access layer

Example repository structure:
```go
package repository

import (
	"net/url"
	"service/internal/pkg/config"
	"service/internal/pkg/core"
	errorpkg "service/internal/pkg/error"
	formpkg "service/internal/pkg/form"
	"service/internal/pkg/model"

	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	"gorm.io/gorm"
)

type CakeRepository interface {
	core.TransactionRepository

	core.PaginateRepository[model.Cake]
	core.FirstIdRepository[model.Cake]

	Store(form formpkg.CakeForm, sellPrice float64) model.Cake
	Delete(cake model.Cake)
	Update(cake model.Cake, form formpkg.CakeForm, sellPrice float64) model.Cake

	AddRecipes(cake model.Cake, recipes []formpkg.CakeCompIngredientForm) []model.CakeRecipeIngredient
	UpdateRecipes(cake model.Cake, recipes []formpkg.CakeCompIngredientForm) []model.CakeRecipeIngredient
	DeleteRecipes(cake model.Cake)

	AddCosts(cake model.Cake, costs []formpkg.CakeCompCostForm) []model.CakeCost
	UpdateCosts(cake model.Cake, costs []formpkg.CakeCompCostForm) []model.CakeCost
	DeleteCosts(cake model.Cake)

	PreloadRecipesAndCosts(query *gorm.DB) *gorm.DB
}

func NewCakeRepository(args ...*gorm.DB) CakeRepository {
	repository := cakeRepository{}
	if len(args) > 0 {
		repository.transaction = args[0]
	} else {
		repository.transaction = config.PgSQL // Default to global config
	}

	return &repository
}

type cakeRepository struct {
	transaction *gorm.DB
}

func (repo *cakeRepository) SetTransaction(tx *gorm.DB) {
	repo.transaction = tx
}

func (repo *cakeRepository) FirstById(id any, args ...func(query *gorm.DB) *gorm.DB) model.Cake {
	var cake model.Cake
	query := repo.transaction

	if len(args) > 0 {
		query = args[0](query)
	}

	err := query.First(&cake, "id = ?", id).Error
	if err != nil {
		errorpkg.ErrXtremeCakeGet(err.Error())
	}

	return cake
}

func (repo *cakeRepository) Paginate(parameter url.Values) ([]model.Cake, interface{}, error) {
	fromDate, toDate := core.SetDateRange(parameter)

	query := repo.transaction.Where("\"createdAt\" BETWEEN ? AND ?", fromDate, toDate)

	if search := parameter.Get("search"); len(search) > 3 {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	cakes, pagination, err := xtrememodel.Paginate(query.Order("id DESC"), parameter, model.Cake{})
	if err != nil {
		return nil, nil, err
	}

	return cakes, pagination, nil
}

func (repo *cakeRepository) Store(form formpkg.CakeForm, sellPrice float64) model.Cake {
	cake := model.Cake{
		Name:        form.Name,
		Description: form.Description,
		Margin:      form.Margin,
		SellPrice:   sellPrice,
		Unit:        form.Unit,
		Stock:       form.Stock,
	}

	err := repo.transaction.Create(&cake).Error
	if err != nil {
		errorpkg.ErrXtremeCakeSave(err.Error())
	}

	return cake
}

func (repo *cakeRepository) Update(cake model.Cake, form formpkg.CakeForm, sellPrice float64) model.Cake {
	cake.Name = form.Name
	cake.Description = form.Description
	cake.Margin = form.Margin
	cake.SellPrice = sellPrice
	cake.Unit = form.Unit
	cake.Stock = form.Stock

	err := repo.transaction.Save(&cake).Error
	if err != nil {
		errorpkg.ErrXtremeCakeSave(err.Error())
	}

	return cake
}

func (repo *cakeRepository) Delete(cake model.Cake) {
	err := repo.transaction.Delete(&cake).Error
	if err != nil {
		errorpkg.ErrXtremeCakeDelete(err.Error())
	}
}

func (repo *cakeRepository) addRecipe(cake model.Cake, recipe formpkg.CakeCompIngredientForm) model.CakeRecipeIngredient {
	cakeRecipe := model.CakeRecipeIngredient{
		CakeID:       cake.ID,
		IngredientID: recipe.IngredientID,
		Amount:       recipe.Amount,
		Unit:         recipe.Unit,
	}

	err := repo.transaction.Create(&cakeRecipe).Error
	if err != nil {
		errorpkg.ErrXtremeCakeRecipeSave(err.Error())
	}

	return cakeRecipe
}

func (repo *cakeRepository) AddRecipes(cake model.Cake, recipes []formpkg.CakeCompIngredientForm) []model.CakeRecipeIngredient {
	var cakeRecipes []model.CakeRecipeIngredient
	for _, recipe := range recipes {
		cakeRecipe := repo.addRecipe(cake, recipe)
		cakeRecipes = append(cakeRecipes, cakeRecipe)
	}
	return cakeRecipes
}

func (repo *cakeRepository) UpdateRecipes(cake model.Cake, recipes []formpkg.CakeCompIngredientForm) []model.CakeRecipeIngredient {
	repo.DeleteRecipes(cake)
	return repo.AddRecipes(cake, recipes)
}

func (repo *cakeRepository) DeleteRecipes(cake model.Cake) {
	err := repo.transaction.Where("\"cakeId\" = ?", cake.ID).Delete(&model.CakeRecipeIngredient{}).Error
	if err != nil {
		errorpkg.ErrXtremeCakeRecipeDelete(err.Error())
	}
}

func (repo *cakeRepository) addCost(cake model.Cake, cost formpkg.CakeCompCostForm) model.CakeCost {
	cakeCost := model.CakeCost{
		CakeID: cake.ID,
		Type:   cost.CostType,
		Cost:   cost.Cost,
	}

	err := repo.transaction.Create(&cakeCost).Error
	if err != nil {
		errorpkg.ErrXtremeCakeCostSave(err.Error())
	}

	return cakeCost
}

func (repo *cakeRepository) AddCosts(cake model.Cake, costs []formpkg.CakeCompCostForm) []model.CakeCost {
	var cakeCosts []model.CakeCost
	for _, cost := range costs {
		cakeCost := repo.addCost(cake, cost)
		cakeCosts = append(cakeCosts, cakeCost)
	}
	return cakeCosts
}

func (repo *cakeRepository) UpdateCosts(cake model.Cake, costs []formpkg.CakeCompCostForm) []model.CakeCost {
	// Delete existing and replace with new costs
	repo.DeleteCosts(cake)
	return repo.AddCosts(cake, costs)
}

func (repo *cakeRepository) DeleteCosts(cake model.Cake) {
	err := repo.transaction.Where("\"cakeId\" = ?", cake.ID).Delete(&model.CakeCost{}).Error
	if err != nil {
		errorpkg.ErrXtremeCakeCostDelete(err.Error())
	}
}

func (repo *cakeRepository) PreloadRecipesAndCosts(query *gorm.DB) *gorm.DB {
	return query.Preload("Recipes.Ingredient").Preload("Costs")
}
```

- Use `service/` for business logic
  - Every service that would modify the data (Create/Update/Delete) should use a transaction and save their action as `Activity`
Example service structure for CRUD:
```go
package service

import (
	"fmt"
	"service/internal/cake/repository"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	"service/internal/pkg/form"
	"service/internal/pkg/model"
	"service/internal/pkg/parser"

	"gorm.io/gorm"
)

type CakeService interface {
	Create(form form.CakeForm) model.Cake
	Update(form form.CakeForm, id string) model.Cake
	Delete(id string) bool
}

func NewCakeService() CakeService {
	return &cakeService{}
}

type cakeService struct {
	repository repository.CakeRepository
}

func (srv *cakeService) Create(form form.CakeForm) model.Cake {
	var cake model.Cake

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)

		cake = srv.repository.Store(form, srv.calculateSellPrice(form))
		recipes := srv.repository.AddRecipes(cake, form.Ingredients)
		costs := srv.repository.AddCosts(cake, form.Costs)

		cake.Recipes = append(cake.Recipes, recipes...)
		cake.Costs = append(cake.Costs, costs...)

		activity.UseActivity{}.SetReference(cake).SetParser(&parser.CakeParser{Object: cake}).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Created new cake: %s [%d]", cake.Name, cake.ID))

		return nil
	})

	return cake
}

func (srv *cakeService) Update(form form.CakeForm, id string) model.Cake {
	var cake model.Cake

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)
		cake = srv.repository.FirstById(id)

		act := activity.UseActivity{}.
			SetReference(cake).
			SetParser(&parser.CakeParser{Object: cake}).
			SetOldProperty(constant.ACTION_UPDATE)

		cake = srv.repository.Update(cake, form, srv.calculateSellPrice(form))
		recipes := srv.repository.UpdateRecipes(cake, form.Ingredients)
		costs := srv.repository.UpdateCosts(cake, form.Costs)

		cake.Recipes = append(cake.Recipes, recipes...)
		cake.Costs = append(cake.Costs, costs...)

		act.SetParser(&parser.CakeParser{Object: cake}).
			SetNewProperty(constant.ACTION_UPDATE).
			Save(fmt.Sprintf("Updated cake: %s [%d]", cake.Name, cake.ID))

		return nil
	})

	return cake
}

func (srv *cakeService) Delete(id string) bool {
	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)
		cake := srv.repository.FirstById(id)

		srv.repository.DeleteRecipes(cake)
		srv.repository.DeleteCosts(cake)
		srv.repository.Delete(cake)

		activity.UseActivity{}.SetReference(cake).SetParser(&parser.CakeParser{Object: cake}).SetOldProperty(constant.ACTION_DELETE).
			Save(fmt.Sprintf("Deleted cake: %s [%d]", cake.Name, cake.ID))

		return nil
	})
	return true
}

func (srv *cakeService) calculateSellPrice(form form.CakeForm) float64 {
	var sellPrice float64

	// Calculate total cost from recipes
	ingredientRepo := repository.NewCakeComponentIngredientRepository()

	var ingredientIDs []any
	recipeQtys := make(map[uint]float64)
	for _, recipe := range form.Ingredients {
		ingredientIDs = append(ingredientIDs, recipe.IngredientID)
		recipeQtys[recipe.IngredientID] = recipe.Amount
	}

	ingredients := ingredientRepo.FindByIds(ingredientIDs)
	for _, ingredient := range ingredients {
		sellPrice += ingredient.UnitPrice * recipeQtys[ingredient.ID]
	}

	// Add additional costs
	for _, cost := range form.Costs {
		sellPrice += cost.Cost
	}

	// Calculate sell price based on margin
	if form.Margin > 0 {
		return sellPrice + (sellPrice * form.Margin / 100)
	}

	return sellPrice
}

```

- (Optional, only append if requested) Use `mail/` for email sending and templating
- (Optional, only append if requested) Use `excel/` for Excel file processing and generation
example excel structure:
```go
package excel

import (
	xtremecore "github.com/globalxtreme/go-core/v2"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/xuri/excelize/v2"
	"service/internal/pkg/constant"
)

type TestingExcel struct {
}

func (ex TestingExcel) Generate() error {
	sheets, properties := ex.setSheetsAndProperties()

	excel := ex.newFile(sheets, properties)
	excel = ex.modifySheet(excel)

	err := excel.Save(constant.PathExcelTesting(), "testing.xlsx")
	if err != nil {
		return err
	}

	return nil
}

func (ex TestingExcel) newFile(sheets []string, properties [][][]interface{}) xtremecore.Excel {
	excel := xtremecore.Excel{
		Sheets:     sheets,
		Properties: properties,
		IsPublic:   false,
	}

	excel.NewFile()

	return excel
}

func (ex TestingExcel) modifySheet(excel xtremecore.Excel) xtremecore.Excel {
	excel.SetStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "860A35", Style: 1},
			{Type: "right", Color: "860A35", Style: 1},
			{Type: "top", Color: "860A35", Style: 1},
			{Type: "bottom", Color: "860A35", Style: 1},
		},
	}, "A1:C1")

	excel.SetStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "A3B763", Style: 2},
			{Type: "right", Color: "A3B763", Style: 2},
			{Type: "top", Color: "A3B763", Style: 2},
			{Type: "bottom", Color: "A3B763", Style: 2},
		},
	}, "A3:C4")

	excel.MergeCells("A5:C5", "A6:A7")

	excel.SetWidthCols([]xtremecore.ColWidth{
		{Cells: "A", Width: 5},
		{Cells: "B", Width: 15},
		{Cells: "C", Width: 5},
	})

	excel.SetHeightRows([]xtremecore.RowHeight{
		{Row: 1, Height: 15},
	})

	return excel
}

func (TestingExcel) setSheetsAndProperties() ([]string, [][][]interface{}) {
	var sheets []string
	var properties [][][]interface{}

	sheets = append(sheets, "Testing 1", "Testing 2")
	for _ = range sheets {
		properties = append(properties, [][]interface{}{
			{"ID", "Name", "Age"},
		})
	}

	dataProperties := [][]int{
		{29, 23, 12},
		{19, 26, 14},
	}

	for sKey, property := range dataProperties {
		for pKey, val := range property {
			properties[sKey] = append(properties[sKey], []interface{}{pKey + 1, xtremepkg.RandomString(10), val})
		}
	}

	return sheets, properties
}
```

- (Optional, only append if requested) Use `pdf/` for PDF generation

### API Route Handlers
- When creating new handlers, *YOU MUST ALWAYS* use cli command: `go run ./cmd/main.go gen:handler <NAME> --type=<web/mobile> --resource`
- Use `--type=web` for web API handlers, `--type=mobile` for mobile API handlers
- Use `--resource` flag to generate resourceful handlers
- Implement handlers in `internal/app/api/<type>/handler/`
- Register created handlers in `internal/app/api/<type>/router.go`

Example handler structure:
```go


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
