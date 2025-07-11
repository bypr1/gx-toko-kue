package model

import (
	"time"

	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

// Transaction stores sales transaction information
// Table: pos_transactions
type Transaction struct {
	xtrememodel.BaseModel
	CakeID          uint      `gorm:"column:cakeId;not null"`
	CakeName        string    `gorm:"column:cakeName;type:varchar(250);not null"`
	Quantity        int       `gorm:"column:quantity;not null"`
	UnitPrice       float64   `gorm:"column:unitPrice;not null"`
	TotalPrice      float64   `gorm:"column:totalPrice;not null"`
	TransactionTime time.Time `gorm:"column:transactionTime;not null"`
	Cake            Cake      `gorm:"foreignKey:CakeID;references:CakeID"`
}

func (Transaction) TableName() string {
	return "pos_transactions"
}

func (t Transaction) SetReference() uint {
	return t.BaseModel.ID
}
