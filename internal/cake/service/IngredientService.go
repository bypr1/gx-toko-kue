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
	"service/internal/pkg/port"

	"gorm.io/gorm"
)

type IngredientService interface {
	SetTransaction(tx *gorm.DB)
	SetActivityRepository(repo port.ActivityRepository)

	Create(form form2.IngredientForm) cake.Ingredient
	Update(form form2.IngredientForm, id uint) cake.Ingredient
	Delete(id uint) error
}

func NewIngredientService() IngredientService {
	return &ingredientService{}
}

type ingredientService struct {
	tx *gorm.DB

	repository         repository.IngredientRepository
	activityRepository port.ActivityRepository
}

func (srv *ingredientService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *ingredientService) SetActivityRepository(repo port.ActivityRepository) {
	srv.activityRepository = repo
}

func (srv *ingredientService) Create(form form2.IngredientForm) cake.Ingredient {
	var ingredient cake.Ingredient

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewIngredientRepository(tx)

		ingredient = srv.repository.Store(form)

		activity.UseActivity{}.SetReference(ingredient).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Created new ingredient: %s [%d]", ingredient.Name, ingredient.ID))

		return nil
	})

	return ingredient
}

func (srv *ingredientService) Update(form form2.IngredientForm, id uint) cake.Ingredient {
	var ingredient cake.Ingredient

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewIngredientRepository(tx)

		ingredient = srv.repository.FirstById(id)
		if ingredient.ID == 0 {
			error2.ErrXtremeIngredientGet("Ingredient not found")
		}

		ingredient = srv.repository.Update(ingredient, form)

		activity.UseActivity{}.SetReference(ingredient).SetNewProperty(constant.ACTION_UPDATE).
			Save(fmt.Sprintf("Updated ingredient: %s [%d]", ingredient.Name, ingredient.ID))

		return nil
	})

	return ingredient
}

func (srv *ingredientService) Delete(id uint) error {
	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewIngredientRepository(tx)
		ingredient := srv.repository.FirstById(id)
		if ingredient.ID == 0 {
			error2.ErrXtremeIngredientGet("Ingredient not found")
		}

		srv.repository.Delete(ingredient)

		activity.UseActivity{}.SetReference(ingredient).SetNewProperty(constant.ACTION_DELETE).
			Save(fmt.Sprintf("Deleted ingredient: %s [%d]", ingredient.Name, ingredient.ID))

		return nil
	})
	return nil
}
