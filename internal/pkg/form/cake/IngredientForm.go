package form

import (
	"net/http"

	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

// IngredientForm is used for validating Ingredient create/update requests
// Maps to model.Ingredient

type IngredientForm struct {
	Name        string  `json:"name" validate:"required,max=250"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price" validate:"required,gte=0"`
	Unit        string  `json:"unit" validate:"required,max=50"`
}

func (f *IngredientForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(f)
}

func (f *IngredientForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &f)
}
