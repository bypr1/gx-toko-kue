package form

import (
	"net/http"

	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type CakeComponentIngredientForm struct {
	Name        string  `json:"name" validate:"required,max=250"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unitPrice" validate:"required,gte=0"`
	UnitId      string  `json:"unitId" validate:"required,max=50"`
}

func (f *CakeComponentIngredientForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(f)
}

func (f *CakeComponentIngredientForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &f)
}
