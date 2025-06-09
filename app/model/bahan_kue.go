package model

type KueBahan struct {
	BahanId int `json:"bahan_id"`
	KueId   int `json:"kue_id"`
	Jumlah  int `json:"jumlah"`
	Bahan   Bahan
}

func (k *KueBahan) TableName() string {
	return "kue_bahan"
}
