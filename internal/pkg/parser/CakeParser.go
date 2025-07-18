package parser

import (
	"service/internal/pkg/constant"
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
	cake := parser.Object
	return map[string]interface{}{
		"id":          cake.ID,
		"name":        cake.Name,
		"description": cake.Description,
		"margin":      cake.Margin,
		"price":       cake.Price,
		"unit":        constant.AllUnitOfMeasure{}.IDAndName(cake.UnitId),
		"stock":       cake.Stock,
		"image":       cake.Image,
		"createdAt":   cake.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":   cake.UpdatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser CakeParser) First() interface{} {
	cake := parser.Object
	var recipes []interface{}
	for _, recipe := range cake.Recipes {
		recipes = append(recipes, CakeRecipeParser{Object: recipe}.First())
	}
	var costs []interface{}
	for _, cost := range cake.Costs {
		costs = append(costs, CakeCostParser{Object: cost}.First())
	}
	return map[string]interface{}{
		"id":          cake.ID,
		"name":        cake.Name,
		"description": cake.Description,
		"margin":      cake.Margin,
		"price":       cake.Price,
		"unit":        constant.AllUnitOfMeasure{}.IDAndName(cake.UnitId),
		"stock":       cake.Stock,
		"createdAt":   cake.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":   cake.UpdatedAt.Format("02/01/2006 15:04"),
		"recipes":     recipes,
		"costs":       costs,
		"image":       cake.Image,
	}
}
