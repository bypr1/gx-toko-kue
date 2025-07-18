package form

import (
	"net/http"

	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type CakeComponentIngredientForm struct {
	Name        string  `json:"name" validate:"required,max=250"`
	Description *string `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	UnitId      int     `json:"unitId" validate:"required"`
}

func (f *CakeComponentIngredientForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(f)
}

func (f *CakeComponentIngredientForm) APIParse(r *http.Request) {
	form := core.BaseForm{}
	form.APIParse(r, &f)
}
