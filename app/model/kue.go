package model

import (
	"gorm.io/gorm"
)

type Kue struct {
	gorm.Model
	ID                   int      `json:"id" gorm:"primaryKey;autoIncrement"`
	Nama                 string   `json:"nama"`
	ProduksiHarian       int      `json:"produksi_harian"`
	HargaTerakhir        int      `json:"harga_terakhir"`
	KeuntunganDiinginkan int      `json:"keuntungan_diinginkan"`
	Bahan                []*Bahan `gorm:"many2many:bahan_kue;"`
}
