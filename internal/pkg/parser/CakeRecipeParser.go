package parser

import (
	"service/internal/pkg/constant"
	"service/internal/pkg/model"
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
		"unit":   constant.UnitOfMeasure{}.IDAndName(recipe.UnitId),
		"ingredient": map[string]interface{}{
			"id":   recipe.IngredientId,
			"name": recipe.Ingredient.Name,
			"unit": constant.UnitOfMeasure{}.IDAndName(recipe.Ingredient.UnitId),
		},
	}
}
