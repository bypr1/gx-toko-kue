package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type CakeIngredient struct {
	xtrememodel.BaseModel
	CakeId       uint    `gorm:"column:cakeId;not null"`
	IngredientId uint    `gorm:"column:ingredientId;not null"`
	Amount       float64 `gorm:"column:amount;not null"` // Amount of ingredient used in the recipe
	UnitId       int     `gorm:"column:unitId;not null"`

	Cake       Cake                    `gorm:"foreignKey:CakeId"`
	Ingredient CakeComponentIngredient `gorm:"foreignKey:IngredientId"`
}

func (CakeIngredient) TableName() string {
	return "cake_ingredients"
}

func (r CakeIngredient) SetReference() uint {
	return r.BaseModel.ID
}
