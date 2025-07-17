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
	PreloadRecipesAndCosts(query *gorm.DB) *gorm.DB

	core.PaginateRepository[model.Cake]
	core.FirstIdRepository[model.Cake]

	Store(form formpkg.CakeForm, sellPrice float64, image string) model.Cake
	Delete(cake model.Cake)
	Update(cake model.Cake, form formpkg.CakeForm, sellPrice float64, image string) model.Cake
	FindByIds(ids []any) []model.Cake

	SaveRecipes(cake model.Cake, recipes []formpkg.CakeFormComponentIngredient) []model.CakeIngredient
	SaveCosts(cake model.Cake, costs []formpkg.CakeFormComponentCost) []model.CakeCost

	DeleteRecipes(cake model.Cake)
	DeleteCosts(cake model.Cake)
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

// -- Preload necessary relations for the model ---

func (repo *cakeRepository) PreloadRecipesAndCosts(query *gorm.DB) *gorm.DB {
	return query.Preload("Recipes.Ingredient").Preload("Costs")
}

// -- Public operations that interact with the database ---

func (repo *cakeRepository) FirstById(id any, args ...func(query *gorm.DB) *gorm.DB) model.Cake {
	var cake model.Cake
	query := config.PgSQL

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

	query := config.PgSQL.Where("\"createdAt\" BETWEEN ? AND ?", fromDate, toDate)

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

	err := config.PgSQL.Where("id IN ?", ids).Find(&cakes).Error
	if err != nil {
		errorpkg.ErrXtremeCakeGet(err.Error())
	}

	return cakes
}

func (repo *cakeRepository) Store(form formpkg.CakeForm, sellPrice float64, image string) model.Cake {
	cake := model.Cake{
		Name:        form.Name,
		Description: form.Description,
		Margin:      form.Margin,
		Price:       sellPrice,
		UnitId:      form.UnitId,
		Stock:       form.Stock,
		Image:       image,
	}

	err := repo.transaction.Create(&cake).Error
	if err != nil {
		errorpkg.ErrXtremeCakeSave(err.Error())
	}

	return cake
}

func (repo *cakeRepository) Update(cake model.Cake, form formpkg.CakeForm, sellPrice float64, image string) model.Cake {
	cake.Name = form.Name
	cake.Description = form.Description
	cake.Margin = form.Margin
	cake.Price = sellPrice
	cake.UnitId = form.UnitId
	cake.Stock = form.Stock
	cake.Image = image

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

func (repo *cakeRepository) SaveRecipes(cake model.Cake, form []formpkg.CakeFormComponentIngredient) []model.CakeIngredient {
	var toDelete []uint
	var toUpdate []model.CakeIngredient
	var toCreate []model.CakeIngredient

	existingRecipes := repo.getExistingRecipes(cake)
	existingMap := make(map[uint]model.CakeIngredient)
	for _, existing := range existingRecipes {
		existingMap[existing.ID] = existing
	}

	for _, f := range form {
		if f.Deleted {
			toDelete = append(toDelete, f.ID)
		} else {
			cakeRecipe := model.CakeIngredient{
				CakeId:       cake.ID,
				IngredientId: f.IngredientId,
				Amount:       f.Amount,
				UnitId:       f.UnitId,
			}

			if f.ID > 0 { // Update existing
				cakeRecipe.ID = f.ID
				toUpdate = append(toUpdate, cakeRecipe)
			} else { // Create new
				toCreate = append(toCreate, cakeRecipe)
			}
		}
	}

	if err := repo.batchDeleteRecipes(toDelete); err != nil {
		errorpkg.ErrXtremeCakeRecipeDelete(err.Error())
	}

	if err := repo.batchUpdateRecipes(toUpdate); err != nil {
		errorpkg.ErrXtremeCakeRecipeSave(err.Error())
	}

	if err := repo.batchCreateRecipes(toCreate); err != nil {
		errorpkg.ErrXtremeCakeRecipeSave(err.Error())
	}

	return append(toUpdate, toCreate...)
}

func (repo *cakeRepository) SaveCosts(cake model.Cake, form []formpkg.CakeFormComponentCost) []model.CakeCost {
	var toDelete []uint
	var toUpdate []model.CakeCost
	var toCreate []model.CakeCost

	existingCosts := repo.getExistingCosts(cake)
	existingMap := make(map[uint]model.CakeCost)
	for _, existing := range existingCosts {
		existingMap[existing.ID] = existing
	}

	for _, f := range form {
		if f.Deleted {
			toDelete = append(toDelete, f.ID)
		} else {
			cakeCost := model.CakeCost{
				CakeId: cake.ID,
				TypeId: f.CostTypeId,
				Cost:   f.Cost,
			}

			if f.ID > 0 { // Update existing
				cakeCost.ID = f.ID
				toUpdate = append(toUpdate, cakeCost)
			} else { // Create new
				toCreate = append(toCreate, cakeCost)
			}
		}
	}

	if err := repo.batchDeleteCosts(toDelete); err != nil {
		errorpkg.ErrXtremeCakeCostDelete(err.Error())
	}

	if err := repo.batchUpdateCosts(toUpdate); err != nil {
		errorpkg.ErrXtremeCakeCostSave(err.Error())
	}

	if err := repo.batchCreateCosts(toCreate); err != nil {
		errorpkg.ErrXtremeCakeCostSave(err.Error())
	}

	return append(toUpdate, toCreate...)
}

func (repo *cakeRepository) DeleteRecipes(cake model.Cake) {
	err := repo.transaction.Where("\"cakeId\" = ?", cake.ID).Delete(&model.CakeIngredient{}).Error
	if err != nil {
		errorpkg.ErrXtremeCakeRecipeDelete(err.Error())
	}
}

func (repo *cakeRepository) DeleteCosts(cake model.Cake) {
	err := repo.transaction.Where("\"cakeId\" = ?", cake.ID).Delete(&model.CakeCost{}).Error
	if err != nil {
		errorpkg.ErrXtremeCakeCostDelete(err.Error())
	}
}

// -- Private helper sections for the repository ---

func (repo *cakeRepository) getExistingRecipes(cake model.Cake) []model.CakeIngredient {
	var recipes []model.CakeIngredient
	repo.transaction.Where("\"cakeId\" = ?", cake.ID).Find(&recipes)
	return recipes
}

func (repo *cakeRepository) batchDeleteRecipes(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	return repo.transaction.Where("\"id\" IN ?", ids).Delete(&model.CakeIngredient{}).Error
}

func (repo *cakeRepository) batchUpdateRecipes(recipes []model.CakeIngredient) error {
	if len(recipes) == 0 {
		return nil
	}

	for _, recipe := range recipes {
		err := repo.transaction.Save(&recipe).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *cakeRepository) batchCreateRecipes(recipes []model.CakeIngredient) error {
	if len(recipes) == 0 {
		return nil
	}

	return repo.transaction.Create(&recipes).Error
}

func (repo *cakeRepository) getExistingCosts(cake model.Cake) []model.CakeCost {
	var costs []model.CakeCost
	repo.transaction.Where("\"cakeId\" = ?", cake.ID).Find(&costs)
	return costs
}

func (repo *cakeRepository) batchDeleteCosts(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	return repo.transaction.Where("\"id\" IN ?", ids).Delete(&model.CakeCost{}).Error
}

func (repo *cakeRepository) batchUpdateCosts(costs []model.CakeCost) error {
	if len(costs) == 0 {
		return nil
	}

	for _, cost := range costs {
		err := repo.transaction.Save(&cost).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *cakeRepository) batchCreateCosts(costs []model.CakeCost) error {
	if len(costs) == 0 {
		return nil
	}

	return repo.transaction.Create(&costs).Error
}
