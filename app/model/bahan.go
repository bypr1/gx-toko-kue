package model

import (
	"gorm.io/gorm"
)

type Bahan struct {
	gorm.Model
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Kue   []*Kue `gorm:"many2many:bahan_kue;"`
}
