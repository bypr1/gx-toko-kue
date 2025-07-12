package port

import (
	"service/internal/pkg/model"
)

/** --- CAKE --- */

type CakeRepository interface {
	FindByIds(ids []any) []model.Cake
}
