package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type CakeComponentIngredient struct {
	xtrememodel.BaseModel
	UnitId      int     `gorm:"column:unitId;not null"`
	Name        string  `gorm:"column:name;type:varchar(250);not null"`
	Description *string `gorm:"column:description;type:text;default:null"`
	Price       float64 `gorm:"column:price;not null"`
}

func (CakeComponentIngredient) TableName() string {
	return "cake_component_ingredients"
}

func (i CakeComponentIngredient) SetReference() uint {
	return i.BaseModel.ID
}
