package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type TransactionCake struct {
	xtrememodel.BaseModel
	TransactionID uint        `gorm:"column:transactionId;not null"`
	CakeID        uint        `gorm:"column:cakeId;not null"`
	Quantity      int         `gorm:"column:quantity;not null"`
	UnitPrice     float64     `gorm:"column:unitPrice;not null"`
	SubTotal      float64     `gorm:"column:subTotal;not null"`
	Transaction   Transaction `gorm:"belongs_to:transactions;foreignKey:TransactionID"`
	Cake          Cake        `gorm:"belongs_to:cakes;foreignKey:CakeID"`
}

func (TransactionCake) TableName() string {
	return "transaction_detail_cakes"
}

func (td TransactionCake) SetReference() uint {
	return td.BaseModel.ID
}
