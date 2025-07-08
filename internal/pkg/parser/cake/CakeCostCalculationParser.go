package cakeparser

import (
	"service/internal/pkg/model/cake"
)

type CakeCostCalculationParser struct {
	Cake      cake.Cake
	TotalCost float64
}

func (parser CakeCostCalculationParser) First() interface{} {
	c := parser.Cake
	totalCost := parser.TotalCost
	return map[string]interface{}{
		"cakeId":    c.ID,
		"cakeName":  c.Name,
		"totalCost": totalCost,
		"margin":    c.Margin,
		"sellPrice": c.SellPrice,
		"profit":    c.SellPrice - totalCost,
	}
}
