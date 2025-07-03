package cake

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

// CakeCost stores production cost information per cake
// Table: cake_costs
type CakeCost struct {
	xtrememodel.BaseModel
	CakeID uint    `gorm:"column:cake_id;not null"`
	Type   string  `gorm:"column:type;type:varchar(100);not null"`
	Cost   float64 `gorm:"column:cost;not null"`
}

func (CakeCost) TableName() string {
	return "cake_costs"
}

func (c CakeCost) SetReference() uint {
	return c.BaseModel.ID
}
