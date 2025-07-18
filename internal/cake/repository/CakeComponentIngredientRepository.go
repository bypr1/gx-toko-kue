package repository

import (
	"net/url"
	"service/internal/pkg/config"
	"service/internal/pkg/core"
	errorpkg "service/internal/pkg/error"
	"service/internal/pkg/form"
	"service/internal/pkg/model"

	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	"gorm.io/gorm"
)

type CakeComponentIngredientRepository interface {
	core.TransactionRepository

	core.PaginateRepository[model.CakeComponentIngredient]
	core.FirstIdRepository[model.CakeComponentIngredient]
	FindByIds(ids []any, args ...func(query *gorm.DB) *gorm.DB) []model.CakeComponentIngredient

	Store(form form.CakeComponentIngredientForm) model.CakeComponentIngredient
	Update(ingredient model.CakeComponentIngredient, form form.CakeComponentIngredientForm) model.CakeComponentIngredient
	Delete(ingredient model.CakeComponentIngredient)
}

func NewCakeComponentIngredientRepository(args ...*gorm.DB) CakeComponentIngredientRepository {
	repository := cakeComponentIngredientRepository{}
	if len(args) > 0 {
		repository.transaction = args[0]
	}

	return &repository
}

type cakeComponentIngredientRepository struct {
	transaction *gorm.DB
}

func (repo *cakeComponentIngredientRepository) SetTransaction(tx *gorm.DB) {
	repo.transaction = tx
}

func (repo *cakeComponentIngredientRepository) FirstById(id any, args ...func(query *gorm.DB) *gorm.DB) model.CakeComponentIngredient {
	var ingredient model.CakeComponentIngredient

	query := config.PgSQL

	if len(args) > 0 {
		query = args[0](query)
	}

	err := query.First(&ingredient, "id = ?", id).Error
	if err != nil {
		errorpkg.ErrXtremeIngredientGet(err.Error())
	}

	return ingredient
}

func (repo *cakeComponentIngredientRepository) Paginate(parameter url.Values) ([]model.CakeComponentIngredient, interface{}, error) {
	fromDate, toDate := core.SetDateRange(parameter)

	query := config.PgSQL.
		Where("\"createdAt\" BETWEEN ? AND ?", fromDate, toDate)

	if search := parameter.Get("search"); len(search) > 3 {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	ingredients, pagination, err := xtrememodel.Paginate(query.Order("id DESC"), parameter, model.CakeComponentIngredient{})
	if err != nil {
		return nil, nil, err
	}

	return ingredients, pagination, nil
}

func (repo *cakeComponentIngredientRepository) FindByIds(ids []any, args ...func(query *gorm.DB) *gorm.DB) []model.CakeComponentIngredient {
	var ingredients []model.CakeComponentIngredient

	query := config.PgSQL
	if len(args) > 0 {
		query = args[0](query)
	}

	err := query.Where("id IN ?", ids).Find(&ingredients).Error
	if err != nil {
		errorpkg.ErrXtremeIngredientGet(err.Error())
	}

	return ingredients
}

func (repo *cakeComponentIngredientRepository) Store(form form.CakeComponentIngredientForm) model.CakeComponentIngredient {
	ingredient := model.CakeComponentIngredient{
		Name:        form.Name,
		Description: form.Description,
		Price:       form.Price,
		UnitId:      form.UnitId,
	}

	err := repo.transaction.Create(&ingredient).Error
	if err != nil {
		errorpkg.ErrXtremeIngredientSave(err.Error())
	}

	return ingredient
}

func (repo *cakeComponentIngredientRepository) Update(ingredient model.CakeComponentIngredient, form form.CakeComponentIngredientForm) model.CakeComponentIngredient {
	ingredient.Name = form.Name
	ingredient.Description = form.Description
	ingredient.Price = form.Price
	ingredient.UnitId = form.UnitId

	err := repo.transaction.Save(&ingredient).Error
	if err != nil {
		errorpkg.ErrXtremeIngredientSave(err.Error())
	}

	return ingredient
}

func (repo *cakeComponentIngredientRepository) Delete(ingredient model.CakeComponentIngredient) {
	err := repo.transaction.Delete(&ingredient).Error
	if err != nil {
		errorpkg.ErrXtremeIngredientDelete(err.Error())
	}
}
