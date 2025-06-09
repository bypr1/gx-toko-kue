package handler

import (
	"github.com/bypr1/gx-toko-kue/app/model"
	"gorm.io/gorm"
	"github.com/labstack/echo/v4"
)

type hargaHandler struct {
	db *gorm.DB
}

func NewHargaHandler(db *gorm.DB) *hargaHandler {
	return &hargaHandler{
		db: db,
	}
}

func (h *hargaHandler) HitungHPPPerKue(e echo.Context) error {
	kueId := e.Param("kue_id")

	//Ambil data kue
	var kue model.Kue
	if err := h.db.First(&kue, kueId).Error; err != nil {
		return e.JSON(500, map[string]interface{}{
			"error": "Gagal mengambil data kue",
		})
	}

	//Hitung total harga komposisi bahan
	var bahan2 []model.KueBahan
	if err := h.db.Find(&bahan2).Preload("Bahan").Where("kue_id = ?", kue.ID).Error; err != nil {
		return e.JSON(500, map[string]interface{}{
			"error": "Gagal mengambil data kue",
		})
	}

	totalHargaBahan := 0
	for _, bahan  := range bahan2 {
		totalHargaBahan = totalHargaBahan + (bahan.Jumlah * bahan.Bahan.Harga)
	}

	//Hitung total harga buruh
	var buruh []model.TenagaKerja
	if err := h.db.Find(&buruh).Error; err != nil {
		return e.JSON(500, map[string]interface{}{
			"error": "Gagal mengambil data buruh",
		})
	}

	for _, bahan  := range buruh {
		totalHargaBuruh = totalHargaBuruh + (bahan.Jumlah * bahan.Bahan.Harga)
	}

	//Hitung berdasarkan upah harian
	totalHargaBuruh += 


	return nil
}
