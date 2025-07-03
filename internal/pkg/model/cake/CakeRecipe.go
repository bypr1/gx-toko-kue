package cake

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

// CakeRecipe stores cake recipe information
// Table: cake_recipes
type CakeRecipe struct {
	xtrememodel.BaseModel
	CakeID       uint    `gorm:"column:cake_id;not null"`
	IngredientID uint    `gorm:"column:ingredient_id;not null"`
	Amount       float64 `gorm:"column:amount;not null"` // Amount of ingredient used in the recipe
	Unit         string  `gorm:"column:unit;type:varchar(50);not null"`
}

func (CakeRecipe) TableName() string {
	return "cake_recipes"
}

func (r CakeRecipe) SetReference() uint {
	return r.BaseModel.ID
}
