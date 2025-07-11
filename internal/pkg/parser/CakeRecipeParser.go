package parser

import "service/internal/pkg/model"

type CakeRecipeParser struct {
	Array  []model.CakeRecipeIngredient
	Object model.CakeRecipeIngredient
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
		"id":           recipe.ID,
		"cakeId":       recipe.CakeID,
		"ingredientId": recipe.IngredientID,
		"amount":       recipe.Amount,
		"unit":         recipe.Unit,
		"ingredient": map[string]interface{}{
			"name": recipe.Ingredient.Name,
			"unit": recipe.Ingredient.Unit,
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
		"unit":         recipe.Unit,
	}
}
