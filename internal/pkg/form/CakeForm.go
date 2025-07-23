package form

import (
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type CakeForm struct {
	Name        string  `json:"name" form:"name" validate:"required,max=250"`
	Description *string `json:"description" form:"description"`
	Margin      float64 `json:"margin" form:"margin" validate:"required,gte=0"`
	UnitId      int     `json:"unitId" form:"unitId" validate:"required"`
	Stock       int     `json:"stock" form:"stock" validate:"gte=0"`

	Ingredients []CakeIngredientForm `json:"ingredients" form:"ingredients" validate:"required,dive"`
	Costs       []CakeCostForm       `json:"costs" form:"costs" validate:"required,dive"`
	Request     *http.Request
}

type CakeIngredientForm struct {
	ID           uint    `json:"id" form:"id"`
	IngredientId uint    `json:"ingredientId" form:"ingredientId" validate:"required_without=ID"`
	Amount       float64 `json:"amount" form:"amount" validate:"required_without=ID"`
	UnitId       int     `json:"unitId" form:"unitId" validate:"required_without=ID"`
	Deleted      bool    `json:"deleted" form:"deleted"`
}

type CakeCostForm struct {
	ID      uint    `json:"id" form:"id"`
	TypeId  int     `json:"typeId" form:"typeId" validate:"required_without=ID"`
	Price   float64 `json:"price" form:"price" validate:"required_without=ID"`
	Deleted bool    `json:"deleted" form:"deleted"`
}

func (f *CakeForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(f)
}

func (f *CakeForm) APIParse(r *http.Request) {
	f.Request = r
	core.BaseForm{}.APIMultipartParse(r, &f)
}
