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

	xtremefs "github.com/globalxtreme/go-core/v2/filesystem"

	"gorm.io/gorm"
)

type CakeService interface {
	Create(form form.CakeForm) model.Cake
	Update(form form.CakeForm, id string) model.Cake
	Delete(id string) bool
}

func NewCakeService() CakeService {
	return &cakeService{}
}

type cakeService struct {
	repository repository.CakeRepository
}

func (srv *cakeService) Create(form form.CakeForm) model.Cake {
	var cake model.Cake
	srv.prepare()

	sellPrice := srv.calculateSellPrice(form)
	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)

		cake = srv.repository.Store(form, sellPrice)
		cake.Recipes = srv.repository.SaveRecipes(cake, form.Ingredients)
		cake.Costs = srv.repository.SaveCosts(cake, form.Costs)

		if image := srv.uploadImage(form); image != "" {
			cake = srv.repository.UpdateImage(cake, image)
		}

		activity.UseActivity{}.SetReference(cake).SetParser(&parser.CakeParser{Object: cake}).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Created new cake: %s [%d]", cake.Name, cake.ID))

		return nil
	})

	return cake
}

func (srv *cakeService) Update(form form.CakeForm, id string) model.Cake {
	cake := srv.prepareWithData(id)

	sellPrice := srv.calculateSellPrice(form)
	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)

		act := activity.UseActivity{}.
			SetReference(cake).
			SetParser(&parser.CakeParser{Object: cake}).
			SetOldProperty(constant.ACTION_UPDATE)

		cake = srv.repository.Update(cake, form, sellPrice)
		cake.Recipes = srv.repository.SaveRecipes(cake, form.Ingredients)
		cake.Costs = srv.repository.SaveCosts(cake, form.Costs)

		if image := srv.uploadImage(form); image != "" {
			cake = srv.repository.UpdateImage(cake, image)
		}

		act.SetParser(&parser.CakeParser{Object: cake}).
			SetNewProperty(constant.ACTION_UPDATE).
			Save(fmt.Sprintf("Updated cake: %s [%d]", cake.Name, cake.ID))

		return nil
	})

	return cake
}

func (srv *cakeService) Delete(id string) bool {
	cake := srv.prepareWithData(id)

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)

		srv.repository.DeleteRecipes(cake)
		srv.repository.DeleteCosts(cake)
		srv.repository.Delete(cake)

		activity.UseActivity{}.SetReference(cake).SetParser(&parser.CakeParser{Object: cake}).SetOldProperty(constant.ACTION_DELETE).
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
		ingredientIDs = append(ingredientIDs, recipe.IngredientId)
		recipeQtys[recipe.IngredientId] = recipe.Amount
	}

	ingredients := ingredientRepo.FindByIds(ingredientIDs)
	for _, ingredient := range ingredients {
		sellPrice += ingredient.Price * recipeQtys[ingredient.ID]
	}

	// Add additional costs
	for _, cost := range form.Costs {
		sellPrice += cost.Price
	}

	// Calculate sell price based on margin
	if form.Margin > 0 {
		return sellPrice + (sellPrice * form.Margin / 100)
	}

	return sellPrice
}

func (srv *cakeService) prepare() {
	srv.repository = repository.NewCakeRepository()
}

func (srv *cakeService) prepareWithData(id any) model.Cake {
	srv.prepare()
	return srv.repository.FirstById(id)
}

func (srv *cakeService) uploadImage(form form.CakeForm) string {
	uploader := xtremefs.Uploader{Path: constant.PathImageCake(), IsPublic: true}
	filePath, err := uploader.MoveFile(form.Request, "image")
	if err != nil {
		return ""
	}

	storage := xtremefs.Storage{IsPublic: uploader.IsPublic}

	return storage.GetFullPathURL(filePath.(string))
}
