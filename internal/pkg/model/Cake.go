package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type Cake struct {
	xtrememodel.BaseModel
	UnitId      int     `gorm:"column:unitId;null"`
	Name        string  `gorm:"column:name;type:varchar(250);not null"`
	Description *string `gorm:"column:description;type:text"`
	Margin      float64 `gorm:"column:margin;not null"`
	Price       float64 `gorm:"column:price;not null"`
	Stock       int     `gorm:"column:stock;default:0"`
	Image       string  `gorm:"column:image;type:varchar(255);default:null"`

	Recipes []CakeIngredient `gorm:"foreignKey:CakeId"`
	Costs   []CakeCost       `gorm:"foreignKey:CakeId"`
}

func (Cake) TableName() string {
	return "cakes"
}

func (c Cake) SetReference() uint {
	return c.BaseModel.ID
}
