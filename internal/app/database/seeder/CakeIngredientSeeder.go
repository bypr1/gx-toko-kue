package seeder

import (
	"service/internal/pkg/model"

	"gorm.io/gorm"
)

type CakeIngredientSeeder struct {
	Connection *gorm.DB
}

func (seed *CakeIngredientSeeder) Seed() {
	ingredients := seed.setIngredientsData()
	for _, ingredient := range ingredients {
		var count int64
		seed.Connection.Model(&model.CakeComponentIngredient{}).Where("name = ?", ingredient["name"]).Count(&count)
		if count > 0 {
			continue
		}

		seed.Connection.Create(&model.CakeComponentIngredient{
			Name:        ingredient["name"].(string),
			Description: ingredient["description"].(string),
			Price:       ingredient["unit_price"].(float64),
			UnitId:      ingredient["unit"].(int),
		})
	}
}

// --- UNEXPORTED FUNCTIONS ---

func (seed *CakeIngredientSeeder) setIngredientsData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "Eggs",
			"description": "Fresh chicken eggs",
			"unit_price":  2000.0,
			"unit":        1,
		},
		{
			"name":        "Flour",
			"description": "All-purpose wheat flour",
			"unit_price":  10000.0,
			"unit":        2,
		},
		{
			"name":        "Sugar",
			"description": "Granulated white sugar",
			"unit_price":  12000.0,
			"unit":        2,
		},
		{
			"name":        "Chocolate",
			"description": "Chocolate compound for baking",
			"unit_price":  25000.0,
			"unit":        2,
		},
	}
}
