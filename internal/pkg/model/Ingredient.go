package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

// Ingredient stores ingredient information
// Table: cake_ingredients
type Ingredient struct {
	xtrememodel.BaseModel
	Name        string  `gorm:"column:name;type:varchar(250);not null"`
	Description string  `gorm:"column:description;type:text;default:null"`
	UnitPrice   float64 `gorm:"column:unitPrice;not null"`
	Unit        string  `gorm:"column:unit;type:varchar(50);not null"`
}

func (Ingredient) TableName() string {
	return "cake_ingredients"
}

func (i Ingredient) SetReference() uint {
	return i.BaseModel.ID
}
