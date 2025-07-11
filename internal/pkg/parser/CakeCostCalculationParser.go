package parser

import "service/internal/pkg/model"

type CakeCostCalculationParser struct {
	Cake      model.Cake
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
