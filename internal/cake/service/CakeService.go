package service

import (
	"fmt"
	"service/internal/cake/repository"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	errorpkg "service/internal/pkg/error"
	"service/internal/pkg/form"
	"service/internal/pkg/model"
	"service/internal/pkg/parser"
	"service/internal/pkg/port"

	"gorm.io/gorm"
)

type CakeService interface {
	SetTransaction(tx *gorm.DB)

	Create(form form.CakeForm) model.Cake
	Update(form form.CakeForm, id any) model.Cake
	Delete(id any) bool
}

func NewCakeService() CakeService {
	return &cakeService{
		tx: config.PgSQL,
	}
}

type cakeService struct {
	tx *gorm.DB

	repository         repository.CakeRepository
	activityRepository port.ActivityRepository
}

func (srv *cakeService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *cakeService) SetActivityRepository(repo port.ActivityRepository) {
	srv.activityRepository = repo
}

func (srv *cakeService) Create(form form.CakeForm) model.Cake {
	var cakeModel model.Cake

	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)

		cakeModel = srv.repository.Store(form, srv.calculateSellPrice(form))
		recipes := srv.repository.AddRecipes(cakeModel, form.Ingredients)
		costs := srv.repository.AddCosts(cakeModel, form.Costs)

		cakeModel.Recipes = append(cakeModel.Recipes, recipes...)
		cakeModel.Costs = append(cakeModel.Costs, costs...)

		activity.InitCreate(cakeModel, tx).
			ParseNewProperty(&parser.CakeParser{Object: cakeModel}).
			Save(fmt.Sprintf("Created new cake: %s [%d]", cakeModel.Name, cakeModel.ID))

		return nil
	})

	return cakeModel
}

func (srv *cakeService) Update(form form.CakeForm, id any) model.Cake {
	var cakeModel model.Cake

	cakeModel = srv.firstOrFail(id)
	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)

		act := activity.InitUpdate(cakeModel, tx).
			ParseOldProperty(&parser.CakeParser{Object: cakeModel})

		cakeModel = srv.repository.Update(cakeModel, form, srv.calculateSellPrice(form))
		recipes := srv.repository.UpdateRecipes(cakeModel, form.Ingredients)
		costs := srv.repository.UpdateCosts(cakeModel, form.Costs)

		cakeModel.Recipes = append(cakeModel.Recipes, recipes...)
		cakeModel.Costs = append(cakeModel.Costs, costs...)

		act.ParseNewProperty(&parser.CakeParser{Object: cakeModel}).
			Save(fmt.Sprintf("Updated cake: %s [%d]", cakeModel.Name, cakeModel.ID))

		return nil
	})

	return cakeModel
}

func (srv *cakeService) Delete(id any) bool {
	cakeModel := srv.firstOrFail(id)
	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)

		srv.repository.DeleteRecipes(cakeModel)
		srv.repository.DeleteCosts(cakeModel)
		srv.repository.Delete(cakeModel)

		activity.InitDelete(cakeModel, tx).
			ParseOldProperty(&parser.CakeParser{Object: cakeModel}).
			Save(fmt.Sprintf("Deleted cake: %s [%d]", cakeModel.Name, cakeModel.ID))

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
	if len(ingredients) == 0 {
		errorpkg.ErrXtremeIngredientGet("No ingredients found for the cake")
	}

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

func (srv *cakeService) firstOrFail(id any) model.Cake {
	cakeModel := srv.repository.FirstById(id)
	if cakeModel.ID == 0 {
		errorpkg.ErrXtremeCakeGet("Cake not found")
	}
	return cakeModel
}
