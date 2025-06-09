package model

import (
	"time"
)

type RiwayatHPP struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	KueId     int       `json:"kue_id"`
	Nama      string    `json:"nama"`
	Harga     int       `json:"harga"`
	CreatedAt time.Time `json:"created_at"`
}
