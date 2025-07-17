package model

import (
	"time"

	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type Transaction struct {
	xtrememodel.BaseModel
	TransactionDate time.Time `gorm:"column:transactionDate;not null"`
	TotalAmount     float64   `gorm:"column:totalAmount;not null"`

	Cakes []TransactionCake `gorm:"foreignKey:TransactionId"`
}

func (Transaction) TableName() string {
	return "transactions"
}

func (t Transaction) SetReference() uint {
	return t.BaseModel.ID
}
