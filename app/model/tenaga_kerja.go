package model

type TenagaKerja struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Nama        string `json:"nama"`
	BiayaHarian int    `json:"biaya_harian"`
}

func (k *TenagaKerja) TableName() string {
	return "tenaga_kerja"
}
