package model

type Bahan struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Kue   []*Kue `gorm:"many2many:bahan_kue;"`
}

func (k *Bahan) TableName() string {
	return "bahan"
}
