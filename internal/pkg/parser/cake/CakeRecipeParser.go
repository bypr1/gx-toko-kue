package cakeparser

import "service/internal/pkg/model/cake"

type CakeRecipeParser struct {
	Array  []cake.CakeRecipe
	Object cake.CakeRecipe
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
