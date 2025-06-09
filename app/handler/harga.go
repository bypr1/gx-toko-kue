package handler

import (
	"github.com/bypr1/gx-toko-kue/app/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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
	for _, bahan := range bahan2 {
		totalHargaBahan = totalHargaBahan + (bahan.Jumlah * bahan.Bahan.Harga)
	}

	//Hitung total harga buruh2
	var buruh2 []model.TenagaKerja
	if err := h.db.Find(&buruh2).Error; err != nil {
		return e.JSON(500, map[string]interface{}{
			"error": "Gagal mengambil data buruh",
		})
	}

	totalHargaBuruh := 0
	for _, buruh := range buruh2 {
		totalHargaBuruh += buruh.BiayaHarian / kue.ProduksiHarian
	}

	//Hitung berdasarkan upah harian
	totalHPP := totalHargaBahan + totalHargaBuruh

	return e.JSON(200, map[string]interface{}{
		"total_hpp": totalHPP,
	})
}
