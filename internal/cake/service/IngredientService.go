package service

import (
	"fmt"
	"service/internal/cake/repository"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	error2 "service/internal/pkg/error"
	form2 "service/internal/pkg/form/cake"
	"service/internal/pkg/model"
	cakeparser "service/internal/pkg/parser"
	"service/internal/pkg/port"

	"gorm.io/gorm"
)

type IngredientService interface {
	SetTransaction(tx *gorm.DB)
	SetActivityRepository(repo port.ActivityRepository)

	Create(form form2.IngredientForm) model.Ingredient
	Update(form form2.IngredientForm, id uint) model.Ingredient
	Delete(id uint) error
}

func NewIngredientService() IngredientService {
	return &ingredientService{
		tx: config.PgSQL,
	}
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

func (srv *ingredientService) Create(form form2.IngredientForm) model.Ingredient {
	var ingredient model.Ingredient

	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewIngredientRepository(tx)

		ingredient = srv.repository.Store(form)

		srv.recordActivity(ingredient, constant.ACTION_CREATE, nil).
			Save(fmt.Sprintf("Created new ingredient: %s [%d]", ingredient.Name, ingredient.ID))

		return nil
	})

	return ingredient
}

func (srv *ingredientService) Update(form form2.IngredientForm, id uint) model.Ingredient {
	var ingredient model.Ingredient

	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewIngredientRepository(tx)

		ingredient = srv.repository.FirstById(id)
		if ingredient.ID == 0 {
			error2.ErrXtremeIngredientGet("Ingredient not found")
		}

		act := srv.recordActivity(ingredient, constant.ACTION_UPDATE, nil)

		ingredient = srv.repository.Update(ingredient, form)

		srv.recordActivity(ingredient, constant.ACTION_UPDATE, &act).
			Save(fmt.Sprintf("Updated ingredient: %s [%d]", ingredient.Name, ingredient.ID))

		return nil
	})

	return ingredient
}

func (srv *ingredientService) Delete(id uint) error {
	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewIngredientRepository(tx)
		ingredient := srv.repository.FirstById(id)
		if ingredient.ID == 0 {
			error2.ErrXtremeIngredientGet("Ingredient not found")
		}

		srv.repository.Delete(ingredient)

		srv.recordActivity(ingredient, constant.ACTION_DELETE, nil).
			Save(fmt.Sprintf("Deleted ingredient: %s [%d]", ingredient.Name, ingredient.ID))

		return nil
	})
	return nil
}

func (srv *ingredientService) recordActivity(ingredient model.Ingredient, action string, act *activity.UseActivity) activity.UseActivity {
	var activ activity.UseActivity
	if act == nil {
		activ = activity.UseActivity{}.
			SetConnection(srv.tx).
			SetReference(&ingredient)

		activ = activ.SetParser(&cakeparser.IngredientParser{Object: ingredient})
		if action != constant.ACTION_CREATE {
			activ = activ.SetOldProperty(action)
		} else {
			activ = activ.SetNewProperty(action)
		}
	} else {
		activ = act.SetReference(&ingredient).
			SetParser(&cakeparser.IngredientParser{Object: ingredient}).
			SetNewProperty(action)
	}

	return activ
}
