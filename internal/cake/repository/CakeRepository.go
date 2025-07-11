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
	FindByIds(ids []any) []model.Cake

	AddRecipes(cake model.Cake, recipes []formpkg.CakeCompIngredientForm) []model.CakeRecipeIngredient
	UpdateRecipes(cake model.Cake, recipes []formpkg.CakeCompIngredientForm) []model.CakeRecipeIngredient
	DeleteRecipes(cake model.Cake)

	AddCosts(cake model.Cake, costs []formpkg.CakeCompCostForm) []model.CakeCost
	UpdateCosts(cake model.Cake, costs []formpkg.CakeCompCostForm) []model.CakeCost
	DeleteCosts(cake model.Cake)

	WithRecipesAndCosts(query *gorm.DB) *gorm.DB
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

func (repo *cakeRepository) FindByIds(ids []any) []model.Cake {
	var cakes []model.Cake

	err := repo.transaction.Where("id IN ?", ids).Find(&cakes).Error
	if err != nil {
		errorpkg.ErrXtremeCakeGet(err.Error())
	}

	return cakes
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

func (repo *cakeRepository) WithRecipesAndCosts(query *gorm.DB) *gorm.DB {
	return query.Preload("Recipes.Ingredient").Preload("Costs")
}
