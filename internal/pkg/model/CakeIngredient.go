package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type CakeIngredient struct {
	xtrememodel.BaseModel
	CakeID       uint    `gorm:"column:cakeId;not null"`
	IngredientID uint    `gorm:"column:ingredientId;not null"`
	Amount       float64 `gorm:"column:amount;not null"` // Amount of ingredient used in the recipe
	Unit         string  `gorm:"column:unit;type:varchar(50);not null"`

	Cake       Cake                    `gorm:"foreignKey:CakeID"`
	Ingredient CakeComponentIngredient `gorm:"foreignKey:IngredientID"`
}

func (CakeIngredient) TableName() string {
	return "cake_recipe_ingredients"
}

func (r CakeIngredient) SetReference() uint {
	return r.BaseModel.ID
}
