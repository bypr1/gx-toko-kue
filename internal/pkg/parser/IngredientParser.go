package parser

import "service/internal/pkg/model"

type IngredientParser struct {
	Array  []model.CakeComponentIngredient
	Object model.CakeComponentIngredient
}

func (parser IngredientParser) CreateActivity(action string) interface{} {
	return parser.First()
}

func (parser IngredientParser) DeleteActivity(action string) interface{} {
	return parser.First()
}

func (parser IngredientParser) GeneralActivity(action string) interface{} {
	return parser.First()
}

func (parser IngredientParser) UpdateActivity(action string) interface{} {
	return parser.First()
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
