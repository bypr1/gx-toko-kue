package service

import (
	"fmt"
	"service/internal/cake/repository"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	error2 "service/internal/pkg/error"
	form2 "service/internal/pkg/form/cake"
	"service/internal/pkg/model"
	cakeparser "service/internal/pkg/parser"
	"service/internal/pkg/port"

	"gorm.io/gorm"
)

type CakeService interface {
	SetTransaction(tx *gorm.DB)

	Create(form form2.CakeForm) model.Cake
	Update(form form2.CakeForm, id uint) model.Cake
	Delete(id uint) error
	CalculateCakeCost(cakeModel model.Cake) float64
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

func (srv *cakeService) Create(form form2.CakeForm) model.Cake {
	var cakeModel model.Cake

	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)

		cakeModel = srv.repository.Store(form)

		for _, recipe := range form.Ingredients {
			cakeRecipe := srv.repository.AddRecipe(cakeModel, recipe)
			cakeModel.Recipes = append(cakeModel.Recipes, cakeRecipe)
		}

		for _, cost := range form.Costs {
			cakeCost := srv.repository.AddCost(cakeModel, cost)
			cakeModel.Costs = append(cakeModel.Costs, cakeCost)
		}

		activity.InitCreate(cakeModel).
			ParseNewProperty(&cakeparser.CakeParser{Object: cakeModel}).
			Save(fmt.Sprintf("Created new cake: %s [%d]", cakeModel.Name, cakeModel.ID))

		return nil
	})

	return cakeModel
}

func (srv *cakeService) Update(form form2.CakeForm, id uint) model.Cake {
	var cakeModel model.Cake

	cakeModel = srv.firstOrFail(id)
	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)

		act := activity.InitUpdate(cakeModel).
			ParseOldProperty(&cakeparser.CakeParser{Object: cakeModel})

		cakeModel = srv.repository.Update(cakeModel, form)

		srv.repository.UpdateRecipes(cakeModel, form.Ingredients)
		srv.repository.UpdateCosts(cakeModel, form.Costs)

		act.ParseNewProperty(&cakeparser.CakeParser{Object: cakeModel}).
			Save(fmt.Sprintf("Updated cake: %s [%d]", cakeModel.Name, cakeModel.ID))

		return nil
	})

	return cakeModel
}

func (srv *cakeService) Delete(id uint) error {
	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)
		cakeModel := srv.repository.FirstById(id)
		if cakeModel.ID == 0 {
			error2.ErrXtremeCakeGet("Cake not found")
		}

		// Delete related recipes and costs first
		srv.repository.DeleteRecipes(cakeModel)
		srv.repository.DeleteCosts(cakeModel)

		// Delete the cake
		srv.repository.Delete(cakeModel)

		activity.InitDelete(cakeModel).
			ParseOldProperty(&cakeparser.CakeParser{Object: cakeModel}).
			Save(fmt.Sprintf("Deleted cake: %s [%d]", cakeModel.Name, cakeModel.ID))

		return nil
	})
	return nil
}

func (srv *cakeService) CalculateCakeCost(cakeModel model.Cake) float64 {
	var totalCost float64

	// Calculate ingredient costs
	ingredientRepo := repository.NewIngredientRepository()
	for _, recipe := range cakeModel.Recipes {
		ingredient := ingredientRepo.FirstById(recipe.IngredientID)
		totalCost += ingredient.UnitPrice * recipe.Amount
	}

	// Add other costs
	for _, cost := range cakeModel.Costs {
		totalCost += cost.Cost
	}

	return totalCost
}

func (srv *cakeService) firstOrFail(id uint) model.Cake {
	cakeModel := srv.repository.FirstById(id)
	if cakeModel.ID == 0 {
		error2.ErrXtremeCakeGet("Cake not found")
	}
	return cakeModel
}
