package parser

import (
	"service/internal/pkg/constant"
	"service/internal/pkg/model"

	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
)

type CakeRecipeParser struct {
	Array  []model.CakeIngredient
	Object model.CakeIngredient
}

func (parser CakeRecipeParser) Get() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, CakeRecipeParser{Object: obj}.First())
	}
	return result
}

func (parser CakeRecipeParser) First() interface{} {
	recipe := parser.Object
	return map[string]interface{}{
		"id":     recipe.ID,
		"amount": recipe.Amount,
		"unit":   constant.CakeUnitOfMeasure{}.Display(xtremepkg.ToInt(recipe.Unit)),
		"ingredient": map[string]interface{}{
			"id":   recipe.IngredientID,
			"name": recipe.Ingredient.Name,
			"unit": constant.CakeIngredientUnitOfMeasure{}.Display(xtremepkg.ToInt(recipe.Ingredient.UnitId)),
		},
	}
}

func (parser CakeRecipeParser) FirstWithoutIngredient() interface{} {
	recipe := parser.Object
	return map[string]interface{}{
		"id":           recipe.ID,
		"cakeId":       recipe.CakeID,
		"ingredientId": recipe.IngredientID,
		"amount":       recipe.Amount,
		"unit":         constant.CakeUnitOfMeasure{}.Display(xtremepkg.ToInt(recipe.Unit)),
	}
}
