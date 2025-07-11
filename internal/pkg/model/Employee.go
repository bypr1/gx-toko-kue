package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

// Employee stores employee information
// Table: hrd_employees
type Employee struct {
	xtrememodel.BaseModel
	Name        string  `gorm:"column:name;type:varchar(250);not null"`
	Position    string  `gorm:"column:position;type:varchar(100);default:null"`
	DailySalary float64 `gorm:"column:dailySalary;not null"`
	DailyHours  float64 `gorm:"column:dailyHours;not null;default:8"`
	IsActive    bool    `gorm:"column:isActive;default:true"`
}

func (Employee) TableName() string {
	return "hrd_employees"
}

func (t Employee) SetReference() uint {
	return t.BaseModel.ID
}
