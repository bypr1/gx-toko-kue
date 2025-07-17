package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type Cake struct {
	xtrememodel.BaseModel
	Name        string  `gorm:"column:name;type:varchar(250);not null"`
	Description string  `gorm:"column:description;type:text;default:null"`
	Margin      float64 `gorm:"column:margin;not null"`
	Price       float64 `gorm:"column:sellPrice;not null"`
	UnitId      int     `gorm:"column:unitId;null"`                          // Unit of measurement for the cake
	Stock       int     `gorm:"column:stock;default:0"`                      // Stock quantity of the cake
	Image       string  `gorm:"column:image;type:varchar(255);default:null"` // Image path for the cake

	Recipes []CakeIngredient `gorm:"foreignKey:CakeId"` // Recipes associated with the cake
	Costs   []CakeCost       `gorm:"foreignKey:CakeId"`
}

func (Cake) TableName() string {
	return "cakes"
}

func (c Cake) SetReference() uint {
	return c.BaseModel.ID
}
