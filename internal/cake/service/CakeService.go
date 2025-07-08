package service

import (
	"fmt"
	"service/internal/cake/repository"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	error2 "service/internal/pkg/error"
	form2 "service/internal/pkg/form/cake"
	"service/internal/pkg/model/cake"
	cakeparser "service/internal/pkg/parser/cake"
	"service/internal/pkg/port"

	"gorm.io/gorm"
)

type CakeService interface {
	SetTransaction(tx *gorm.DB)
	SetActivityRepository(repo port.ActivityRepository)

	Create(form form2.CakeForm) cake.Cake
	Update(form form2.CakeForm, id uint) cake.Cake
	Delete(id uint) error
	CalculateCakeCost(cakeModel cake.Cake) float64
}

func NewCakeService() CakeService {
	return &cakeService{
		tx: config.CakeSQL,
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

func (srv *cakeService) Create(form form2.CakeForm) cake.Cake {
	var cakeModel cake.Cake

	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)

		cakeModel = srv.repository.Store(form)

		// Add recipes
		for _, recipe := range form.Ingredients {
			cakeRecipe := srv.repository.AddRecipe(cakeModel, recipe)
			cakeModel.Recipes = append(cakeModel.Recipes, cakeRecipe)
		}

		// Add costs
		for _, cost := range form.Costs {
			cakeCost := srv.repository.AddCost(cakeModel, cost)
			cakeModel.Costs = append(cakeModel.Costs, cakeCost)
		}

		// Pass the transaction to recordActivity
		srv.recordActivity(cakeModel, constant.ACTION_CREATE, nil).
			Save(fmt.Sprintf("Created new cake: %s [%d]", cakeModel.Name, cakeModel.ID))

		return nil
	})

	return cakeModel
}

func (srv *cakeService) Update(form form2.CakeForm, id uint) cake.Cake {
	var cakeModel cake.Cake

	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)

		cakeModel = srv.repository.FirstById(id)
		if cakeModel.ID == 0 {
			error2.ErrXtremeCakeGet("Cake not found")
		}

		act := srv.recordActivity(cakeModel, constant.ACTION_UPDATE, nil)

		cakeModel = srv.repository.Update(cakeModel, form)

		// Update recipes
		srv.repository.UpdateRecipes(cakeModel, form.Ingredients)

		// Update costs
		srv.repository.UpdateCosts(cakeModel, form.Costs)

		// Reload the cake with updated relationships
		cakeModel = srv.repository.FirstById(id, func(query *gorm.DB) *gorm.DB {
			return query.Preload("Recipes").Preload("Costs")
		})

		srv.recordActivity(cakeModel, constant.ACTION_UPDATE, &act).
			Save(fmt.Sprintf("Updated cake: %s [%d]", cakeModel.Name, cakeModel.ID))

		return nil
	})

	return cakeModel
}

func (srv *cakeService) Delete(id uint) error {
	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeRepository(tx)
		cakeModel := srv.repository.FirstById(id, func(query *gorm.DB) *gorm.DB {
			return query.Preload("Recipes").Preload("Costs")
		})
		if cakeModel.ID == 0 {
			error2.ErrXtremeCakeGet("Cake not found")
		}

		// Delete related recipes and costs first
		srv.repository.DeleteRecipes(cakeModel)
		srv.repository.DeleteCosts(cakeModel)

		// Delete the cake
		srv.repository.Delete(cakeModel)

		srv.recordActivity(cakeModel, constant.ACTION_DELETE, nil).
			Save(fmt.Sprintf("Deleted cake: %s [%d]", cakeModel.Name, cakeModel.ID))

		return nil
	})
	return nil
}

func (srv *cakeService) CalculateCakeCost(cakeModel cake.Cake) float64 {
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

func (srv *cakeService) recordActivity(cake cake.Cake, action string, act *activity.UseActivity) activity.UseActivity {
	var activ activity.UseActivity
	if act == nil {
		activ = activity.UseActivity{}.
			SetConnection(srv.tx).
			SetReference(&cake)

		activ = activ.SetParser(&cakeparser.CakeParser{Object: cake})
		if action != constant.ACTION_CREATE {
			activ = activ.SetOldProperty(action)
		} else {
			activ = activ.SetNewProperty(action)
		}
	} else {
		activ = act.SetReference(&cake).
			SetParser(&cakeparser.CakeParser{Object: cake}).
			SetNewProperty(action)
	}

	return activ
}
