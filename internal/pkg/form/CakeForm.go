package form

import (
	"fmt"
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
)

type CakeForm struct {
	Name        string                   `json:"name" validate:"required,max=250"`
	Description string                   `json:"description"`
	Margin      float64                  `json:"margin" validate:"required,gte=0"`
	Unit        string                   `json:"unit" validate:"max=50"`
	Stock       int                      `json:"stock" validate:"gte=0"`
	Ingredients []CakeCompIngredientForm `json:"ingredients" validate:"required,dive"`
	Costs       []CakeCompCostForm       `json:"costs" validate:"required,dive"`
	Request     *http.Request
}

type CakeCompIngredientForm struct {
	IngredientID uint    `json:"ingredientId" validate:"required,gt=0"`
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

func (f *CakeForm) FormParse(r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		f.APIParse(r)
	}

	f.Name = r.FormValue("name")
	f.Description = r.FormValue("description")
	f.Margin = xtremepkg.ToFloat64(r.FormValue("margin"))
	f.Unit = r.FormValue("unit")
	f.Stock = xtremepkg.ToInt(r.FormValue("stock"))
	f.Ingredients = make([]CakeCompIngredientForm, 0)

	exists := true
	for i := 0; exists; i++ {
		var compIngredient CakeCompIngredientForm
		compIngredient.IngredientID = uint(xtremepkg.ToInt(r.FormValue(fmt.Sprintf("ingredients[%d][ingredientId]", i))))
		if compIngredient.IngredientID == 0 {
			exists = false
			continue
		}

		compIngredient.Amount = xtremepkg.ToFloat64(r.FormValue(fmt.Sprintf("ingredients[%d][amount]", i)))
		compIngredient.Unit = r.FormValue(fmt.Sprintf("ingredients[%d][unit]", i))
		f.Ingredients = append(f.Ingredients, compIngredient)
	}

	exists = true
	f.Costs = make([]CakeCompCostForm, 0)
	for i := 0; exists; i++ {
		var compCost CakeCompCostForm
		compCost.CostType = r.FormValue(fmt.Sprintf("costs[%d][type]", i))
		if compCost.CostType == "" {
			exists = false
			continue
		}

		compCost.Cost = xtremepkg.ToFloat64(r.FormValue(fmt.Sprintf("costs[%d][cost]", i)))
		f.Costs = append(f.Costs, compCost)
	}
	f.Request = r
}
