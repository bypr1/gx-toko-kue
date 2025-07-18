package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type TransactionCake struct {
	xtrememodel.BaseModel
	TransactionId uint    `gorm:"column:transactionId;not null"`
	CakeId        uint    `gorm:"column:cakeId;not null"`
	Quantity      int     `gorm:"column:quantity;not null"`
	Price         float64 `gorm:"column:price;not null"`
	SubTotal      float64 `gorm:"column:subTotal;not null"`

	Transaction Transaction `gorm:"belongs_to:transactions;foreignKey:TransactionId"`
	Cake        Cake        `gorm:"belongs_to:cakes;foreignKey:CakeId"`
}

func (TransactionCake) TableName() string {
	return "transaction_cakes"
}

func (td TransactionCake) SetReference() uint {
	return td.BaseModel.ID
}
