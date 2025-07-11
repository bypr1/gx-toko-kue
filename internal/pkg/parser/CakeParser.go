package parser

import (
	"service/internal/pkg/model"
)

type CakeParser struct {
	Array  []model.Cake
	Object model.Cake
}

func (parser CakeParser) CreateActivity(action string) interface{} {
	return parser.Brief()
}

func (parser CakeParser) DeleteActivity(action string) interface{} {
	return parser.Brief()
}

func (parser CakeParser) GeneralActivity(action string) interface{} {
	return parser.Brief()
}

func (parser CakeParser) UpdateActivity(action string) interface{} {
	return parser.Brief()
}

func (parser CakeParser) Briefs() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, CakeParser{Object: obj}.Brief())
	}
	return result
}

func (parser CakeParser) Brief() interface{} {
	cakeObj := parser.Object
	var recipes []interface{}
	for _, recipe := range cakeObj.Recipes {
		recipes = append(recipes, CakeRecipeParser{Object: recipe}.First())
	}
	var costs []interface{}
	for _, cost := range cakeObj.Costs {
		costs = append(costs, CakeCostParser{Object: cost}.First())
	}
	return map[string]interface{}{
		"id":          cakeObj.ID,
		"name":        cakeObj.Name,
		"description": cakeObj.Description,
		"margin":      cakeObj.Margin,
		"sellPrice":   cakeObj.SellPrice,
		"unit":        cakeObj.Unit,
		"stock":       cakeObj.Stock,
		"createdAt":   cakeObj.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":   cakeObj.UpdatedAt.Format("02/01/2006 15:04"),
		"recipes":     recipes,
		"costs":       costs,
	}
}
