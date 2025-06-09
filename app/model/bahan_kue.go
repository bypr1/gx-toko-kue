package model

import "github.com/jinzhu/gorm"

type KueBahan struct {
	gorm.Model
	BahanId int `json:"bahan_id"`
	KueId   int `json:"kue_id"`
	Jumlah  int `json:"jumlah"`
	Bahan   Bahan
}
