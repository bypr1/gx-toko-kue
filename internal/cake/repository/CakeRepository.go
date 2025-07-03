package repository

import (
	"net/url"
	"service/internal/pkg/config"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	form2 "service/internal/pkg/form/cake"
	"service/internal/pkg/model/cake"

	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	"gorm.io/gorm"
)

type CakeRepository interface {
	core.TransactionRepository

	FirstById(id any, args ...func(query *gorm.DB) *gorm.DB) cake.Cake
	Find(parameter url.Values) ([]cake.Cake, interface{}, error)

	Store(form form2.CakeForm) cake.Cake
	Update(cakeModel cake.Cake, form form2.CakeForm) cake.Cake
	Delete(cakeModel cake.Cake)

	AddRecipe(cakeModel cake.Cake, recipe form2.CakeCompIngredientForm) cake.CakeRecipe
	UpdateRecipes(cakeModel cake.Cake, recipes []form2.CakeCompIngredientForm)
	DeleteRecipes(cakeModel cake.Cake)

	AddCost(cakeModel cake.Cake, cost form2.CakeCompCostForm) cake.CakeCost
	UpdateCosts(cakeModel cake.Cake, costs []form2.CakeCompCostForm)
	DeleteCosts(cakeModel cake.Cake)
}

func NewCakeRepository(args ...*gorm.DB) CakeRepository {
	repository := cakeRepository{}
	if len(args) > 0 {
		repository.transaction = args[0]
	}

	return &repository
}

type cakeRepository struct {
	transaction *gorm.DB
}

func (repo *cakeRepository) SetTransaction(tx *gorm.DB) {
	repo.transaction = tx
}

func (repo *cakeRepository) FirstById(id any, args ...func(query *gorm.DB) *gorm.DB) cake.Cake {
	var cakeModel cake.Cake

	query := config.PgSQL
	if repo.transaction != nil {
		query = repo.transaction
	}

	if len(args) > 0 {
		query = args[0](query)
	}

	err := query.First(&cakeModel, "id = ?", id).Error
	if err != nil {
		error2.ErrXtremeCakeGet(err.Error())
	}

	return cakeModel
}

func (repo *cakeRepository) Find(parameter url.Values) ([]cake.Cake, interface{}, error) {
	fromDate, toDate := core.SetDateRange(parameter)

	query := config.PgSQL.Preload("Recipes").Preload("Costs").
		Where("\"createdAt\" BETWEEN ? AND ?", fromDate, toDate)

	if search := parameter.Get("search"); len(search) > 3 {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	cakes, pagination, err := xtrememodel.Paginate(query.Order("id DESC"), parameter, cake.Cake{})
	if err != nil {
		return nil, nil, err
	}

	return cakes, pagination, nil
}

func (repo *cakeRepository) Store(form form2.CakeForm) cake.Cake {
	db := config.PgSQL
	if repo.transaction != nil {
		db = repo.transaction
	}

	cakeModel := cake.Cake{
		Name:        form.Name,
		Description: form.Description,
		Margin:      form.Margin,
		SellPrice:   form.SellPrice,
		Unit:        form.Unit,
		Stock:       form.Stock,
	}

	err := db.Create(&cakeModel).Error
	if err != nil {
		error2.ErrXtremeCakeSave(err.Error())
	}

	return cakeModel
}

func (repo *cakeRepository) Update(cakeModel cake.Cake, form form2.CakeForm) cake.Cake {
	db := config.PgSQL
	if repo.transaction != nil {
		db = repo.transaction
	}

	cakeModel.Name = form.Name
	cakeModel.Description = form.Description
	cakeModel.Margin = form.Margin
	cakeModel.SellPrice = form.SellPrice
	cakeModel.Unit = form.Unit
	cakeModel.Stock = form.Stock

	err := db.Save(&cakeModel).Error
	if err != nil {
		error2.ErrXtremeCakeSave(err.Error())
	}

	return cakeModel
}

func (repo *cakeRepository) Delete(cakeModel cake.Cake) {
	db := config.PgSQL
	if repo.transaction != nil {
		db = repo.transaction
	}

	err := db.Delete(&cakeModel).Error
	if err != nil {
		error2.ErrXtremeCakeDelete(err.Error())
	}
}

func (repo *cakeRepository) AddRecipe(cakeModel cake.Cake, recipe form2.CakeCompIngredientForm) cake.CakeRecipe {
	db := config.PgSQL
	if repo.transaction != nil {
		db = repo.transaction
	}

	cakeRecipe := cake.CakeRecipe{
		CakeID:       cakeModel.ID,
		IngredientID: recipe.IngredientID,
		Amount:       recipe.Amount,
		Unit:         recipe.Unit,
	}

	err := db.Create(&cakeRecipe).Error
	if err != nil {
		error2.ErrXtremeCakeRecipeSave(err.Error())
	}

	return cakeRecipe
}

func (repo *cakeRepository) UpdateRecipes(cakeModel cake.Cake, recipes []form2.CakeCompIngredientForm) {
	// Delete existing recipes
	repo.DeleteRecipes(cakeModel)

	// Add new recipes
	for _, recipe := range recipes {
		repo.AddRecipe(cakeModel, recipe)
	}
}

func (repo *cakeRepository) DeleteRecipes(cakeModel cake.Cake) {
	db := config.PgSQL
	if repo.transaction != nil {
		db = repo.transaction
	}

	err := db.Where("cake_id = ?", cakeModel.ID).Delete(&cake.CakeRecipe{}).Error
	if err != nil {
		error2.ErrXtremeCakeRecipeDelete(err.Error())
	}
}

func (repo *cakeRepository) AddCost(cakeModel cake.Cake, cost form2.CakeCompCostForm) cake.CakeCost {
	db := config.PgSQL
	if repo.transaction != nil {
		db = repo.transaction
	}

	cakeCost := cake.CakeCost{
		CakeID: cakeModel.ID,
		Type:   cost.CostType,
		Cost:   cost.Cost,
	}

	err := db.Create(&cakeCost).Error
	if err != nil {
		error2.ErrXtremeCakeCostSave(err.Error())
	}

	return cakeCost
}

func (repo *cakeRepository) UpdateCosts(cakeModel cake.Cake, costs []form2.CakeCompCostForm) {
	// Delete existing costs
	repo.DeleteCosts(cakeModel)

	// Add new costs
	for _, cost := range costs {
		repo.AddCost(cakeModel, cost)
	}
}

func (repo *cakeRepository) DeleteCosts(cakeModel cake.Cake) {
	db := config.PgSQL
	if repo.transaction != nil {
		db = repo.transaction
	}

	err := db.Where("cake_id = ?", cakeModel.ID).Delete(&cake.CakeCost{}).Error
	if err != nil {
		error2.ErrXtremeCakeCostDelete(err.Error())
	}
}
