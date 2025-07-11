package service

import (
	"fmt"
	"service/internal/cake/repository"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	errorPkg "service/internal/pkg/error"
	"service/internal/pkg/form"
	"service/internal/pkg/model"
	cakeparser "service/internal/pkg/parser"

	"gorm.io/gorm"
)

type CakeComponentIngredientService interface {
	SetTransaction(tx *gorm.DB)

	Create(form form.CakeComponentIngredientForm) model.CakeComponentIngredient
	Update(form form.CakeComponentIngredientForm, id any) model.CakeComponentIngredient
	Delete(id any) bool
}

func NewIngredientService() CakeComponentIngredientService {
	return &cakeComponentIngredientService{
		tx: config.PgSQL,
	}
}

type cakeComponentIngredientService struct {
	tx *gorm.DB

	repository repository.CakeComponentIngredientRepository
}

func (srv *cakeComponentIngredientService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *cakeComponentIngredientService) Create(form form.CakeComponentIngredientForm) model.CakeComponentIngredient {
	var ingredient model.CakeComponentIngredient

	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeComponentIngredientRepository(tx)

		ingredient = srv.repository.Store(form)

		activity.InitCreate(ingredient, tx).
			ParseNewProperty(&cakeparser.IngredientParser{Object: ingredient}).
			Save(fmt.Sprintf("Created new ingredient: %s [%d]", ingredient.Name, ingredient.ID))

		return nil
	})

	return ingredient
}

func (srv *cakeComponentIngredientService) Update(form form.CakeComponentIngredientForm, id any) model.CakeComponentIngredient {
	var ingredient model.CakeComponentIngredient

	ingredient = srv.firstOrFail(id)

	srv.tx.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCakeComponentIngredientRepository(tx)

		act := activity.InitUpdate(ingredient, tx).
			ParseOldProperty(&cakeparser.IngredientParser{Object: ingredient})

		ingredient = srv.repository.Update(ingredient, form)

		act.ParseNewProperty(&cakeparser.IngredientParser{Object: ingredient}).
			Save(fmt.Sprintf("Updated ingredient: %s [%d]", ingredient.Name, ingredient.ID))

		return nil
	})

	return ingredient
}

func (srv *cakeComponentIngredientService) Delete(id any) bool {
	ingredient := srv.firstOrFail(id)
	srv.tx.Transaction(func(tx *gorm.DB) error {
		repository.NewCakeComponentIngredientRepository(tx).
			Delete(ingredient)

		activity.InitDelete(ingredient, tx).
			ParseOldProperty(&cakeparser.IngredientParser{Object: ingredient}).
			Save(fmt.Sprintf("Deleted ingredient: %s [%d]", ingredient.Name, ingredient.ID))

		return nil
	})
	return true
}

func (srv *cakeComponentIngredientService) firstOrFail(id any) model.CakeComponentIngredient {
	ingredientModel := srv.repository.FirstById(id)
	if ingredientModel.ID == 0 {
		errorPkg.ErrXtremeIngredientGet("Ingredient not found")
	}
	return ingredientModel
}
