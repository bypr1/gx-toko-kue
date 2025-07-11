package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type CakeCost struct {
	xtrememodel.BaseModel
	CakeID uint    `gorm:"column:cakeId;not null"`
	Type   string  `gorm:"column:type;type:varchar(100);not null"`
	Cost   float64 `gorm:"column:cost;not null"`
}

func (CakeCost) TableName() string {
	return "cake_costs"
}

func (c CakeCost) SetReference() uint {
	return c.BaseModel.ID
}
