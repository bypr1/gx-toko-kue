package parser

import (
	"service/internal/pkg/constant"
	"service/internal/pkg/model"
)

type CakeParser struct {
	Array  []model.Cake
	Object model.Cake
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
		"unit":        constant.UnitOfMeasure{}.IDAndName(cake.UnitId),
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
		recipes = append(recipes, CakeIngredientParser{Object: recipe}.First())
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
		"unit":        constant.UnitOfMeasure{}.IDAndName(cake.UnitId),
		"stock":       cake.Stock,
		"createdAt":   cake.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":   cake.UpdatedAt.Format("02/01/2006 15:04"),
		"recipes":     recipes,
		"costs":       costs,
		"image":       cake.Image,
	}
}

type CakeCostParser struct {
	Array  []model.CakeCost
	Object model.CakeCost
}

func (parser CakeCostParser) Get() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, CakeCostParser{Object: obj}.First())
	}
	return result
}

func (parser CakeCostParser) First() interface{} {
	cost := parser.Object
	return map[string]interface{}{
		"id":    cost.ID,
		"type":  constant.CakeCostType{}.IDAndName(cost.TypeId),
		"price": cost.Price,
	}
}

type CakeIngredientParser struct {
	Array  []model.CakeIngredient
	Object model.CakeIngredient
}

func (parser CakeIngredientParser) Get() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, CakeIngredientParser{Object: obj}.First())
	}
	return result
}

func (parser CakeIngredientParser) First() interface{} {
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

func (parser CakeParser) CreateActivity(action string) interface{} {
	return parser.First()
}

func (parser CakeParser) UpdateActivity(action string) interface{} {
	return parser.First()
}

func (parser CakeParser) DeleteActivity(action string) interface{} {
	return parser.First()
}

func (parser CakeParser) GeneralActivity(action string) interface{} {
	return parser.First()
}
