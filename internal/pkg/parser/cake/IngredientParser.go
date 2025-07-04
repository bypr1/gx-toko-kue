package cakeparser

import "service/internal/pkg/model/cake"

type IngredientParser struct {
	Array  []cake.Ingredient
	Object cake.Ingredient
}

func (parser IngredientParser) Get() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, IngredientParser{Object: obj}.First())
	}
	return result
}

func (parser IngredientParser) First() interface{} {
	ingredient := parser.Object
	return map[string]interface{}{
		"id":          ingredient.ID,
		"name":        ingredient.Name,
		"description": ingredient.Description,
		"unitPrice":   ingredient.UnitPrice,
		"unit":        ingredient.Unit,
	}
}
