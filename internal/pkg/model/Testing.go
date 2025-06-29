package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type Testing struct {
	xtrememodel.BaseModel
	Name string `gorm:"column:name;type:varchar(250);default:null"`

	Subs []TestingSub `gorm:"foreignKey:testingId"`
}

func (Testing) TableName() string {
	return "testing"
}

func (model Testing) SetReference() uint {
	return model.BaseModel.ID
}
