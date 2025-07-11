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
	Delete(cakeModel model.Cake)
	Update(cakeModel model.Cake, form formpkg.CakeForm, sellPrice float64) model.Cake

	AddRecipes(cakeModel model.Cake, recipes []formpkg.CakeCompIngredientForm) []model.CakeRecipeIngredient
	UpdateRecipes(cakeModel model.Cake, recipes []formpkg.CakeCompIngredientForm) []model.CakeRecipeIngredient
	DeleteRecipes(cakeModel model.Cake)

	AddCosts(cakeModel model.Cake, costs []formpkg.CakeCompCostForm) []model.CakeCost
	UpdateCosts(cakeModel model.Cake, costs []formpkg.CakeCompCostForm) []model.CakeCost
	DeleteCosts(cakeModel model.Cake)
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
	var cakeModel model.Cake
	query := repo.transaction

	if len(args) > 0 {
		query = args[0](query)
	}

	err := query.First(&cakeModel, "id = ?", id).Error
	if err != nil {
		errorpkg.ErrXtremeCakeGet(err.Error())
	}

	return cakeModel
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
	cakeModel := model.Cake{
		Name:        form.Name,
		Description: form.Description,
		Margin:      form.Margin,
		SellPrice:   sellPrice,
		Unit:        form.Unit,
		Stock:       form.Stock,
	}

	err := repo.transaction.Create(&cakeModel).Error
	if err != nil {
		errorpkg.ErrXtremeCakeSave(err.Error())
	}

	return cakeModel
}

func (repo *cakeRepository) Update(cakeModel model.Cake, form formpkg.CakeForm, sellPrice float64) model.Cake {
	cakeModel.Name = form.Name
	cakeModel.Description = form.Description
	cakeModel.Margin = form.Margin
	cakeModel.SellPrice = sellPrice
	cakeModel.Unit = form.Unit
	cakeModel.Stock = form.Stock

	err := repo.transaction.Save(&cakeModel).Error
	if err != nil {
		errorpkg.ErrXtremeCakeSave(err.Error())
	}

	return cakeModel
}

func (repo *cakeRepository) Delete(cakeModel model.Cake) {
	err := repo.transaction.Delete(&cakeModel).Error
	if err != nil {
		errorpkg.ErrXtremeCakeDelete(err.Error())
	}
}

func (repo *cakeRepository) addRecipe(cakeModel model.Cake, recipe formpkg.CakeCompIngredientForm) model.CakeRecipeIngredient {
	cakeRecipe := model.CakeRecipeIngredient{
		CakeID:       cakeModel.ID,
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

func (repo *cakeRepository) AddRecipes(cakeModel model.Cake, recipes []formpkg.CakeCompIngredientForm) []model.CakeRecipeIngredient {
	var cakeRecipes []model.CakeRecipeIngredient
	for _, recipe := range recipes {
		cakeRecipe := repo.addRecipe(cakeModel, recipe)
		cakeRecipes = append(cakeRecipes, cakeRecipe)
	}
	return cakeRecipes
}

func (repo *cakeRepository) UpdateRecipes(cakeModel model.Cake, recipes []formpkg.CakeCompIngredientForm) []model.CakeRecipeIngredient {
	repo.DeleteRecipes(cakeModel)
	return repo.AddRecipes(cakeModel, recipes)
}

func (repo *cakeRepository) DeleteRecipes(cakeModel model.Cake) {
	err := repo.transaction.Where("\"cakeId\" = ?", cakeModel.ID).Delete(&model.CakeRecipeIngredient{}).Error
	if err != nil {
		errorpkg.ErrXtremeCakeRecipeDelete(err.Error())
	}
}

func (repo *cakeRepository) addCost(cakeModel model.Cake, cost formpkg.CakeCompCostForm) model.CakeCost {
	cakeCost := model.CakeCost{
		CakeID: cakeModel.ID,
		Type:   cost.CostType,
		Cost:   cost.Cost,
	}

	err := repo.transaction.Create(&cakeCost).Error
	if err != nil {
		errorpkg.ErrXtremeCakeCostSave(err.Error())
	}

	return cakeCost
}

func (repo *cakeRepository) AddCosts(cakeModel model.Cake, costs []formpkg.CakeCompCostForm) []model.CakeCost {
	var cakeCosts []model.CakeCost
	for _, cost := range costs {
		cakeCost := repo.addCost(cakeModel, cost)
		cakeCosts = append(cakeCosts, cakeCost)
	}
	return cakeCosts
}

func (repo *cakeRepository) UpdateCosts(cakeModel model.Cake, costs []formpkg.CakeCompCostForm) []model.CakeCost {
	// Delete existing and replace with new costs
	repo.DeleteCosts(cakeModel)
	return repo.AddCosts(cakeModel, costs)
}

func (repo *cakeRepository) DeleteCosts(cakeModel model.Cake) {
	err := repo.transaction.Where("\"cakeId\" = ?", cakeModel.ID).Delete(&model.CakeCost{}).Error
	if err != nil {
		errorpkg.ErrXtremeCakeCostDelete(err.Error())
	}
}
