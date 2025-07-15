package port

import (
	"service/internal/pkg/model"

	"gorm.io/gorm"
)

/** --- CAKE --- */

type CakeRepository interface {
	FindByIds(ids []any) []model.Cake
	SetTransaction(tx *gorm.DB)
}
