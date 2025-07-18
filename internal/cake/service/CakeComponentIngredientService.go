package service

import (
	"fmt"
	"service/internal/cake/repository"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	"service/internal/pkg/form"
	"service/internal/pkg/model"
	cakeparser "service/internal/pkg/parser"

	"gorm.io/gorm"
)

type CakeComponentIngredientService interface {
	Create(form form.CakeComponentIngredientForm) model.CakeComponentIngredient
	Update(form form.CakeComponentIngredientForm, id any) model.CakeComponentIngredient
	Delete(id any) bool
}

func NewIngredientService() CakeComponentIngredientService {
	return &cakeComponentIngredientService{}
}

type cakeComponentIngredientService struct {
	repository repository.CakeComponentIngredientRepository
}

func (srv *cakeComponentIngredientService) Create(form form.CakeComponentIngredientForm) model.CakeComponentIngredient {
	var ingredient model.CakeComponentIngredient
	srv.prepare()

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)

		ingredient = srv.repository.Store(form)

		activity.UseActivity{}.SetReference(ingredient).SetParser(&cakeparser.IngredientParser{Object: ingredient}).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Created new ingredient: %s [%d]", ingredient.Name, ingredient.ID))

		return nil
	})

	return ingredient
}

func (srv *cakeComponentIngredientService) Update(form form.CakeComponentIngredientForm, id any) model.CakeComponentIngredient {
	ingredient := srv.prepareWithData(id)

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)

		act := activity.UseActivity{}.SetReference(ingredient).SetParser(&cakeparser.IngredientParser{Object: ingredient}).SetOldProperty(constant.ACTION_UPDATE)

		ingredient = srv.repository.Update(ingredient, form)

		act.SetParser(&cakeparser.IngredientParser{Object: ingredient}).SetNewProperty(constant.ACTION_UPDATE).
			Save(fmt.Sprintf("Updated ingredient: %s [%d]", ingredient.Name, ingredient.ID))

		return nil
	})

	return ingredient
}

func (srv *cakeComponentIngredientService) Delete(id any) bool {
	ingredient := srv.prepareWithData(id)

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)
		srv.repository.Delete(ingredient)

		activity.UseActivity{}.SetReference(ingredient).SetParser(&cakeparser.IngredientParser{Object: ingredient}).SetOldProperty(constant.ACTION_DELETE).
			Save(fmt.Sprintf("Deleted ingredient: %s [%d]", ingredient.Name, ingredient.ID))

		return nil
	})
	return true
}

func (srv *cakeComponentIngredientService) prepare() {
	srv.repository = repository.NewCakeComponentIngredientRepository(nil)
}

func (srv *cakeComponentIngredientService) prepareWithData(id any) model.CakeComponentIngredient {
	srv.prepare()
	return srv.repository.FirstById(id)
}
