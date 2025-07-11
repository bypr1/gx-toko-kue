package repository

import (
	"net/url"
	"service/internal/pkg/config"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	form2 "service/internal/pkg/form/cake"
	"service/internal/pkg/model"

	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	"gorm.io/gorm"
)

type IngredientRepository interface {
	core.TransactionRepository

	FirstById(id any, args ...func(query *gorm.DB) *gorm.DB) model.Ingredient
	Find(parameter url.Values) ([]model.Ingredient, interface{}, error)

	Store(form form2.IngredientForm) model.Ingredient
	Update(ingredient model.Ingredient, form form2.IngredientForm) model.Ingredient
	Delete(ingredient model.Ingredient)
}

func NewIngredientRepository(args ...*gorm.DB) IngredientRepository {
	repository := ingredientRepository{
		transaction: config.PgSQL, // Default to global config
	}
	if len(args) > 0 {
		repository.transaction = args[0]
	}

	return &repository
}

type ingredientRepository struct {
	transaction *gorm.DB
}

func (repo *ingredientRepository) SetTransaction(tx *gorm.DB) {
	repo.transaction = tx
}

func (repo *ingredientRepository) FirstById(id any, args ...func(query *gorm.DB) *gorm.DB) model.Ingredient {
	var ingredient model.Ingredient

	query := repo.transaction

	if len(args) > 0 {
		query = args[0](query)
	}

	err := query.First(&ingredient, "id = ?", id).Error
	if err != nil {
		error2.ErrXtremeIngredientGet(err.Error())
	}

	return ingredient
}

func (repo *ingredientRepository) Find(parameter url.Values) ([]model.Ingredient, interface{}, error) {
	fromDate, toDate := core.SetDateRange(parameter)

	query := repo.transaction.
		Where("\"createdAt\" BETWEEN ? AND ?", fromDate, toDate)

	if search := parameter.Get("search"); len(search) > 3 {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	ingredients, pagination, err := xtrememodel.Paginate(query.Order("id DESC"), parameter, model.Ingredient{})
	if err != nil {
		return nil, nil, err
	}

	return ingredients, pagination, nil
}

func (repo *ingredientRepository) Store(form form2.IngredientForm) model.Ingredient {
	ingredient := model.Ingredient{
		Name:        form.Name,
		Description: form.Description,
		UnitPrice:   form.UnitPrice,
		Unit:        form.Unit,
	}

	err := repo.transaction.Create(&ingredient).Error
	if err != nil {
		error2.ErrXtremeIngredientSave(err.Error())
	}

	return ingredient
}

func (repo *ingredientRepository) Update(ingredient model.Ingredient, form form2.IngredientForm) model.Ingredient {
	ingredient.Name = form.Name
	ingredient.Description = form.Description
	ingredient.UnitPrice = form.UnitPrice
	ingredient.Unit = form.Unit

	err := repo.transaction.Save(&ingredient).Error
	if err != nil {
		error2.ErrXtremeIngredientSave(err.Error())
	}

	return ingredient
}

func (repo *ingredientRepository) Delete(ingredient model.Ingredient) {
	err := repo.transaction.Delete(&ingredient).Error
	if err != nil {
		error2.ErrXtremeIngredientDelete(err.Error())
	}
}
