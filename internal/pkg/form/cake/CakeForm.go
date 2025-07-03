package form

import (
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type CakeForm struct {
	Name        string                   `json:"name" validate:"required,max=250"`
	Description string                   `json:"description"`
	Margin      float64                  `json:"margin" validate:"required,gte=0"`
	SellPrice   float64                  `json:"sell_price" validate:"required,gte=0"`
	Unit        string                   `json:"unit" validate:"max=50"`
	Stock       int                      `json:"stock" validate:"gte=0"`
	Ingredients []CakeCompIngredientForm `json:"ingredients" validate:"required,dive"`
	Costs       []CakeCompCostForm       `json:"costs" validate:"required,dive"`
}

type CakeCompIngredientForm struct {
	IngredientID uint    `json:"ingredient_id" validate:"required,gt=0"`
	Amount       float64 `json:"amount" validate:"required,gte=0"`
	Unit         string  `json:"unit" validate:"required,max=50"`
}

type CakeCompCostForm struct {
	CostType string  `json:"type" validate:"required,max=100"`
	Cost     float64 `json:"cost" validate:"required,gte=0"`
}

func (f *CakeForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(f)
}

func (f *CakeForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &f)
}
