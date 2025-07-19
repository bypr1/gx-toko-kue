package constant

import (
	"service/internal/pkg/core"
)

const (
	CAKE_COST_TYPE_LABOR          = "labor"
	CAKE_COST_TYPE_LABOR_ID       = 1
	CAKE_COST_TYPE_PACKAGING      = "packaging"
	CAKE_COST_TYPE_PACKAGING_ID   = 2
	CAKE_COST_TYPE_GAS            = "gas"
	CAKE_COST_TYPE_GAS_ID         = 3
	CAKE_COST_TYPE_ELECTRICITY    = "electricity"
	CAKE_COST_TYPE_ELECTRICITY_ID = 4
)

type CakeCostType struct{}

func (c CakeCostType) OptionIDNames() map[int]string {
	return map[int]string{
		CAKE_COST_TYPE_LABOR_ID:       CAKE_COST_TYPE_LABOR,
		CAKE_COST_TYPE_PACKAGING_ID:   CAKE_COST_TYPE_PACKAGING,
		CAKE_COST_TYPE_GAS_ID:         CAKE_COST_TYPE_GAS,
		CAKE_COST_TYPE_ELECTRICITY_ID: CAKE_COST_TYPE_ELECTRICITY,
	}
}

func (c CakeCostType) IDAndName(id int) map[string]interface{} {
	return core.IDName{}.IDAndName(id, c)
}

func (c CakeCostType) Display(id int) string {
	return core.IDName{}.Display(id, c)
}
