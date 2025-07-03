package cakeparser

import "service/internal/pkg/model/cake"

type CakeCostParser struct {
	Array  []cake.CakeCost
	Object cake.CakeCost
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
		"id":     cost.ID,
		"cakeId": cost.CakeID,
		"type":   cost.Type,
		"cost":   cost.Cost,
	}
}
