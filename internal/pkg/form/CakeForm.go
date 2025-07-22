package form

import (
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type CakeForm struct {
	Name        string  `json:"name" validate:"required,max=250"`
	Description *string `json:"description"`
	Margin      float64 `json:"margin" validate:"required,gte=0"`
	UnitId      int     `json:"unitId" validate:"required"`
	Stock       int     `json:"stock" validate:"gte=0"`

	Ingredients []CakeIngredientForm `json:"ingredients" validate:"required,dive"`
	Costs       []CakeCostForm       `json:"costs" validate:"required,dive"`
	Request     *http.Request
}

type CakeIngredientForm struct {
	ID           uint    `json:"id"`
	IngredientId uint    `json:"ingredientId" validate:"required_without=ID"`
	Amount       float64 `json:"amount" validate:"required_without=ID"`
	UnitId       int     `json:"unitId" validate:"required_without=ID"`
	Deleted      bool    `json:"deleted"`
}

type CakeCostForm struct {
	ID      uint    `json:"id"`
	TypeId  int     `json:"typeId" validate:"required_without=ID"`
	Price   float64 `json:"price" validate:"required_without=ID"`
	Deleted bool    `json:"deleted"`
}

func (f *CakeForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(f)
}

func (f *CakeForm) APIParse(r *http.Request) {
	f.Request = r
	core.BaseForm{}.APIMultipartParse(r, &f)
}
