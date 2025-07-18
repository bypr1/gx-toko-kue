package parser

import (
	"service/internal/pkg/constant"
	"service/internal/pkg/model"
)

type CakeComponentIngredientParser struct {
	Array  []model.CakeComponentIngredient
	Object model.CakeComponentIngredient
}

func (parser CakeComponentIngredientParser) Get() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, CakeComponentIngredientParser{Object: obj}.First())
	}
	return result
}

func (parser CakeComponentIngredientParser) First() interface{} {
	ingredient := parser.Object
	return map[string]interface{}{
		"id":          ingredient.ID,
		"name":        ingredient.Name,
		"description": ingredient.Description,
		"price":       ingredient.Price,
		"unit":        constant.UnitOfMeasure{}.IDAndName(ingredient.UnitId),
	}
}

func (parser CakeComponentIngredientParser) CreateActivity(action string) interface{} {
	return parser.First()
}

func (parser CakeComponentIngredientParser) UpdateActivity(action string) interface{} {
	return parser.First()
}

func (parser CakeComponentIngredientParser) DeleteActivity(action string) interface{} {
	return parser.First()
}

func (parser CakeComponentIngredientParser) GeneralActivity(action string) interface{} {
	return parser.First()
}
