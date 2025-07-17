package form

import (
	"fmt"
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
)

type CakeForm struct {
	Name        string                        `json:"name" validate:"required,max=250"`
	Description string                        `json:"description"`
	Margin      float64                       `json:"margin" validate:"required,gte=0"`
	UnitId      int                           `json:"unitId" validate:"max=50"`
	Stock       int                           `json:"stock" validate:"gte=0"`
	Ingredients []CakeFormComponentIngredient `json:"ingredients" validate:"required,dive"`
	Costs       []CakeFormComponentCost       `json:"costs" validate:"required,dive"`
	Request     *http.Request
}

type CakeFormComponentIngredient struct {
	ID           uint    `json:"id"`
	IngredientId uint    `json:"ingredientId" validate:"required,gt=0"`
	Amount       float64 `json:"amount" validate:"required,gte=0"`
	UnitId       int     `json:"unitId" validate:"required,max=50"`
	Deleted      bool    `json:"deleted"`
}

type CakeFormComponentCost struct {
	ID         uint    `json:"id"`
	CostTypeId int     `json:"typeId" validate:"required,max=100"`
	Cost       float64 `json:"cost" validate:"required,gte=0"`
	Deleted    bool    `json:"deleted"`
}

func (f *CakeForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(f)
}

func (f *CakeForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &f)
}

func (f *CakeForm) FormParse(r *http.Request) {
	formValue := r.MultipartForm.Value

	f.Name = formValue["name"][0]
	f.Description = formValue["description"][0]
	f.Margin = xtremepkg.ToFloat64(formValue["margin"][0])
	f.UnitId = xtremepkg.ToInt(formValue["unitId"][0])
	f.Stock = xtremepkg.ToInt(formValue["stock"][0])

	f.Ingredients = make([]CakeFormComponentIngredient, 0)
	isLooping := true
	for i := 0; isLooping; i++ {
		ingredientIDKey := fmt.Sprintf("ingredients[%d][ingredientId]", i)
		if ingredientIDValues, exists := formValue[ingredientIDKey]; exists {
			var compIngredient CakeFormComponentIngredient
			compIngredient.IngredientId = uint(xtremepkg.ToInt(ingredientIDValues[0]))
			compIngredient.Amount = xtremepkg.ToFloat64(formValue[fmt.Sprintf("ingredients[%d][amount]", i)][0])
			compIngredient.UnitId = xtremepkg.ToInt(formValue[fmt.Sprintf("ingredients[%d][unitId]", i)][0])

			idKey := fmt.Sprintf("ingredients[%d][id]", i)
			if idValues, exists := formValue[idKey]; exists && len(idValues) > 0 {
				compIngredient.ID = uint(xtremepkg.ToInt(idValues[0]))
				compIngredient.Deleted = xtremepkg.ToBool(formValue[fmt.Sprintf("ingredients[%d][deleted]", i)][0])
			}

			f.Ingredients = append(f.Ingredients, compIngredient)
		} else {
			isLooping = false
		}
	}

	f.Costs = make([]CakeFormComponentCost, 0)
	isLooping = true
	for i := 0; isLooping; i++ {
		costTypeKey := fmt.Sprintf("costs[%d][typeId]", i)
		if costTypeValues, exists := formValue[costTypeKey]; exists {
			var compCost CakeFormComponentCost
			compCost.CostTypeId = xtremepkg.ToInt(costTypeValues[0])
			compCost.Cost = xtremepkg.ToFloat64(formValue[fmt.Sprintf("costs[%d][cost]", i)][0])

			idKey := fmt.Sprintf("costs[%d][id]", i)
			if idValues, exists := formValue[idKey]; exists && len(idValues) > 0 {
				compCost.ID = uint(xtremepkg.ToInt(idValues[0]))
				compCost.Deleted = xtremepkg.ToBool(formValue[fmt.Sprintf("costs[%d][deleted]", i)][0])
			}

			f.Costs = append(f.Costs, compCost)
		} else {
			isLooping = false
		}
	}
	f.Request = r
}
