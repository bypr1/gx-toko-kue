package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

// CakeComponentIngredient stores ingredient information
// Table: cake_ingredients
type CakeComponentIngredient struct {
	xtrememodel.BaseModel
	Name        string  `gorm:"column:name;type:varchar(250);not null"`
	Description string  `gorm:"column:description;type:text;default:null"`
	UnitPrice   float64 `gorm:"column:unitPrice;not null"`
	Unit        string  `gorm:"column:unit;type:varchar(50);not null"`
}

func (CakeComponentIngredient) TableName() string {
	return "cake_component_ingredients"
}

func (i CakeComponentIngredient) SetReference() uint {
	return i.BaseModel.ID
}
