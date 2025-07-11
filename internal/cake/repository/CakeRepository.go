package repository

import (
	"net/url"
	"service/internal/pkg/config"
	"service/internal/pkg/core"
	gxError "service/internal/pkg/error"
	cakeForm "service/internal/pkg/form/cake"
	"service/internal/pkg/model"

	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	"gorm.io/gorm"
)

type CakeRepository interface {
	core.TransactionRepository

	core.PaginateRepository[model.Cake]
	core.FirstIdRepository[model.Cake]

	Store(form cakeForm.CakeForm) model.Cake
	Delete(cakeModel model.Cake)
	Update(cakeModel model.Cake, form cakeForm.CakeForm) model.Cake

	AddRecipe(cakeModel model.Cake, recipe cakeForm.CakeCompIngredientForm) model.CakeRecipe
	UpdateRecipes(cakeModel model.Cake, recipes []cakeForm.CakeCompIngredientForm)
	DeleteRecipes(cakeModel model.Cake)

	AddCost(cakeModel model.Cake, cost cakeForm.CakeCompCostForm) model.CakeCost
	UpdateCosts(cakeModel model.Cake, costs []cakeForm.CakeCompCostForm)
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
		gxError.ErrXtremeCakeGet(err.Error())
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

func (repo *cakeRepository) Store(form cakeForm.CakeForm) model.Cake {
	cakeModel := model.Cake{
		Name:        form.Name,
		Description: form.Description,
		Margin:      form.Margin,
		SellPrice:   form.SellPrice,
		Unit:        form.Unit,
		Stock:       form.Stock,
	}

	err := repo.transaction.Create(&cakeModel).Error
	if err != nil {
		gxError.ErrXtremeCakeSave(err.Error())
	}

	return cakeModel
}

func (repo *cakeRepository) Update(cakeModel model.Cake, form cakeForm.CakeForm) model.Cake {
	cakeModel.Name = form.Name
	cakeModel.Description = form.Description
	cakeModel.Margin = form.Margin
	cakeModel.SellPrice = form.SellPrice
	cakeModel.Unit = form.Unit
	cakeModel.Stock = form.Stock

	err := repo.transaction.Save(&cakeModel).Error
	if err != nil {
		gxError.ErrXtremeCakeSave(err.Error())
	}

	return cakeModel
}

func (repo *cakeRepository) Delete(cakeModel model.Cake) {
	err := repo.transaction.Delete(&cakeModel).Error
	if err != nil {
		gxError.ErrXtremeCakeDelete(err.Error())
	}
}

func (repo *cakeRepository) AddRecipe(cakeModel model.Cake, recipe cakeForm.CakeCompIngredientForm) model.CakeRecipe {
	cakeRecipe := model.CakeRecipe{
		CakeID:       cakeModel.ID,
		IngredientID: recipe.IngredientID,
		Amount:       recipe.Amount,
		Unit:         recipe.Unit,
	}

	err := repo.transaction.Create(&cakeRecipe).Error
	if err != nil {
		gxError.ErrXtremeCakeRecipeSave(err.Error())
	}

	return cakeRecipe
}

func (repo *cakeRepository) UpdateRecipes(cakeModel model.Cake, recipes []cakeForm.CakeCompIngredientForm) {
	// Delete existing and replace with new recipes
	repo.DeleteRecipes(cakeModel)
	for _, recipe := range recipes {
		repo.AddRecipe(cakeModel, recipe)
	}
}

func (repo *cakeRepository) DeleteRecipes(cakeModel model.Cake) {
	err := repo.transaction.Where("\"cakeId\" = ?", cakeModel.ID).Delete(&model.CakeRecipe{}).Error
	if err != nil {
		gxError.ErrXtremeCakeRecipeDelete(err.Error())
	}
}

func (repo *cakeRepository) AddCost(cakeModel model.Cake, cost cakeForm.CakeCompCostForm) model.CakeCost {
	cakeCost := model.CakeCost{
		CakeID: cakeModel.ID,
		Type:   cost.CostType,
		Cost:   cost.Cost,
	}

	err := repo.transaction.Create(&cakeCost).Error
	if err != nil {
		gxError.ErrXtremeCakeCostSave(err.Error())
	}

	return cakeCost
}

func (repo *cakeRepository) UpdateCosts(cakeModel model.Cake, costs []cakeForm.CakeCompCostForm) {
	// Delete existing and replace with new costs
	repo.DeleteCosts(cakeModel)
	for _, cost := range costs {
		repo.AddCost(cakeModel, cost)
	}
}

func (repo *cakeRepository) DeleteCosts(cakeModel model.Cake) {
	err := repo.transaction.Where("\"cakeId\" = ?", cakeModel.ID).Delete(&model.CakeCost{}).Error
	if err != nil {
		gxError.ErrXtremeCakeCostDelete(err.Error())
	}
}
