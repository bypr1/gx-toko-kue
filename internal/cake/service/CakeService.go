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
	"service/internal/pkg/port"

	"gorm.io/gorm"
)

type CakeService interface {
	SetTransaction(tx *gorm.DB)

	Create(form form.CakeForm) model.Cake
	Update(form form.CakeForm, id string) model.Cake
	Delete(id string) bool
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
	var cake model.Cake

	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)

		cake = srv.repository.Store(form, srv.calculateSellPrice(form))
		recipes := srv.repository.AddRecipes(cake, form.Ingredients)
		costs := srv.repository.AddCosts(cake, form.Costs)

		cake.Recipes = append(cake.Recipes, recipes...)
		cake.Costs = append(cake.Costs, costs...)

		activity.UseActivity{}.SetParser(&parser.CakeParser{Object: cake}).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Created new cake: %s [%d]", cake.Name, cake.ID))

		return nil
	})

	return cake
}

func (srv *cakeService) Update(form form.CakeForm, id string) model.Cake {
	var cake model.Cake

	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)
		cake = srv.repository.FirstById(id)

		act := activity.UseActivity{}.SetParser(&parser.CakeParser{Object: cake}).SetOldProperty(constant.ACTION_UPDATE)

		cake = srv.repository.Update(cake, form, srv.calculateSellPrice(form))
		recipes := srv.repository.UpdateRecipes(cake, form.Ingredients)
		costs := srv.repository.UpdateCosts(cake, form.Costs)

		cake.Recipes = append(cake.Recipes, recipes...)
		cake.Costs = append(cake.Costs, costs...)

		act.SetParser(&parser.CakeParser{Object: cake}).SetNewProperty(constant.ACTION_UPDATE).
			Save(fmt.Sprintf("Updated cake: %s [%d]", cake.Name, cake.ID))

		return nil
	})

	return cake
}

func (srv *cakeService) Delete(id string) bool {
	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)
		cake := srv.repository.FirstById(id)

		srv.repository.DeleteRecipes(cake)
		srv.repository.DeleteCosts(cake)
		srv.repository.Delete(cake)

		activity.UseActivity{}.SetParser(&parser.CakeParser{Object: cake}).SetOldProperty(constant.ACTION_DELETE).
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
