package model

import (
	"gorm.io/gorm"
)

type TenagaKerja struct {
	gorm.Model
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Nama        string `json:"nama"`
	BiayaHarian int    `json:"biaya_harian"`
}
