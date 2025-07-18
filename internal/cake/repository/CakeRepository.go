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

	Store(form formpkg.CakeForm, sellPrice float64) model.Cake
	Delete(cake model.Cake)
	Update(cake model.Cake, form formpkg.CakeForm, sellPrice float64) model.Cake
	UpdateImage(cake model.Cake, image string) model.Cake
	FindByIds(ids []any) []model.Cake

	SaveRecipes(cake model.Cake, recipes []formpkg.CakeIngredientForm) []model.CakeIngredient
	SaveCosts(cake model.Cake, costs []formpkg.CakeCostForm) []model.CakeCost

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

func (repo *cakeRepository) PreloadRecipesAndCosts(query *gorm.DB) *gorm.DB {
	return query.Preload("Recipes.Ingredient").Preload("Costs")
}

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

func (repo *cakeRepository) Store(form formpkg.CakeForm, sellPrice float64) model.Cake {
	cake := model.Cake{
		Name:        form.Name,
		Description: form.Description,
		Margin:      form.Margin,
		Price:       sellPrice,
		UnitId:      form.UnitId,
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
	cake.Price = sellPrice
	cake.UnitId = form.UnitId
	cake.Stock = form.Stock

	err := repo.transaction.Save(&cake).Error
	if err != nil {
		errorpkg.ErrXtremeCakeSave(err.Error())
	}

	return cake
}

func (repo *cakeRepository) UpdateImage(cake model.Cake, image string) model.Cake {
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

func (repo *cakeRepository) SaveRecipes(cake model.Cake, requests []formpkg.CakeIngredientForm) []model.CakeIngredient {
	var recipes []model.CakeIngredient

	for _, request := range requests {
		var recipe model.CakeIngredient
		if request.ID > 0 {
			if request.Deleted {
				err := repo.transaction.Where("id = ?", request.ID).Delete(&model.CakeIngredient{}).Error
				if err != nil {
					errorpkg.ErrXtremeCakeRecipeDelete(err.Error())
				}
			} else {
				repo.transaction.Preload("Ingredient").First(&recipe, "id = ?", request.ID)
				if recipe.ID == 0 {
					errorpkg.ErrXtremeCakeRecipeGet("ID not found")
				}

				recipe.IngredientId = request.IngredientId
				recipe.Amount = request.Amount
				recipe.UnitId = request.UnitId
				err := repo.transaction.Save(&recipe).Error
				if err != nil {
					errorpkg.ErrXtremeCakeRecipeSave(err.Error())
				}
			}
		} else {
			recipe = model.CakeIngredient{
				CakeId:       cake.ID,
				IngredientId: request.IngredientId,
				Amount:       request.Amount,
				UnitId:       request.UnitId,
			}
			err := repo.transaction.Save(&recipe).Error
			repo.transaction.Preload("Ingredient").First(&recipe, "id = ?", recipe.ID)
			if err != nil {
				errorpkg.ErrXtremeCakeRecipeSave(err.Error())
			}
		}

		if recipe.ID > 0 {
			recipes = append(recipes, recipe)
		}
	}

	return recipes
}

func (repo *cakeRepository) SaveCosts(cake model.Cake, requests []formpkg.CakeCostForm) []model.CakeCost {
	var costs []model.CakeCost

	for _, request := range requests {
		var cost model.CakeCost
		if request.ID > 0 {
			if request.Deleted {
				err := repo.transaction.Where("id = ?", request.ID).Delete(&model.CakeCost{}).Error
				if err != nil {
					errorpkg.ErrXtremeCakeCostDelete(err.Error())
				}
			} else {
				repo.transaction.First(&cost, "id = ?", request.ID)
				if cost.ID == 0 {
					errorpkg.ErrXtremeCakeCostGet("ID not found")
				}

				cost.TypeId = request.TypeId
				cost.Price = request.Price
				err := repo.transaction.Save(&cost).Error
				if err != nil {
					errorpkg.ErrXtremeCakeCostSave(err.Error())
				}
			}
		} else {
			cost = model.CakeCost{
				CakeId: cake.ID,
				TypeId: request.TypeId,
				Price:  request.Price,
			}
			err := repo.transaction.Save(&cost).Error
			if err != nil {
				errorpkg.ErrXtremeCakeCostSave(err.Error())
			}
		}

		if cost.ID > 0 {
			costs = append(costs, cost)
		}
	}

	return costs
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
