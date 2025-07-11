package parser

import "service/internal/pkg/model"

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
		"id":   cost.ID,
		"type": cost.Type,
		"cost": cost.Cost,
	}
}
