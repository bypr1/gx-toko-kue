package cakeparser

import (
	"service/internal/pkg/model/cake"
)

type CakeParser struct {
	Array  []cake.Cake
	Object cake.Cake
}

func (parser CakeParser) CreateActivity(action string) interface{} {
	return parser.First()
}

func (parser CakeParser) DeleteActivity(action string) interface{} {
	return parser.First()
}

func (parser CakeParser) GeneralActivity(action string) interface{} {
	return parser.First()
}

func (parser CakeParser) UpdateActivity(action string) interface{} {
	return parser.First()
}

func (parser CakeParser) Get() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, CakeParser{Object: obj}.First())
	}
	return result
}

func (parser CakeParser) First() interface{} {
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
