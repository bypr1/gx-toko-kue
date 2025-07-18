package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type CakeCost struct {
	xtrememodel.BaseModel
	CakeId uint    `gorm:"column:cakeId;not null"`
	TypeId int     `gorm:"column:typeId;not null"`
	Price  float64 `gorm:"column:price;not null"`
}

func (CakeCost) TableName() string {
	return "cake_costs"
}

func (c CakeCost) SetReference() uint {
	return c.BaseModel.ID
}
