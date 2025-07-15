package model

import (
	"time"

	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type Transaction struct {
	xtrememodel.BaseModel
	TransactionDate time.Time               `gorm:"column:transactionDate;not null"`
	TotalAmount     float64                 `gorm:"column:totalAmount;not null"`
	Cakes           []TransactionDetailCake `gorm:"one2many:transaction_details;foreignKey:TransactionID"`
}

func (Transaction) TableName() string {
	return "transactions"
}

func (t Transaction) SetReference() uint {
	return t.BaseModel.ID
}

type TransactionReport struct {
	Transaction
	TotalCakes int `gorm:"column:totalCakes;not null"`
}
