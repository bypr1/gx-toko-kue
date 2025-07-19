package form

import (
	"fmt"
	"net/http"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
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
	formValue := r.MultipartForm.Value

	f.Name = formValue["name"][0]
	f.Description = &formValue["description"][0]
	f.Margin = xtremepkg.ToFloat64(formValue["margin"][0])
	f.UnitId = xtremepkg.ToInt(formValue["unitId"][0])
	f.Stock = xtremepkg.ToInt(formValue["stock"][0])

	f.Ingredients = make([]CakeIngredientForm, 0)
	isLooping := true
	for i := 0; isLooping; i++ {
		var compIngredient CakeIngredientForm
		childOK, idOK := false, false

		ingredientIDKey := fmt.Sprintf("ingredients[%d][ingredientId]", i)
		if ingredientIDValues, exists := formValue[ingredientIDKey]; exists {
			compIngredient.IngredientId = uint(xtremepkg.ToInt(ingredientIDValues[0]))
			compIngredient.Amount = xtremepkg.ToFloat64(formValue[fmt.Sprintf("ingredients[%d][amount]", i)][0])
			compIngredient.UnitId = xtremepkg.ToInt(formValue[fmt.Sprintf("ingredients[%d][unitId]", i)][0])
			childOK = true
		}

		idKey := fmt.Sprintf("ingredients[%d][id]", i)
		if idValues, exists := formValue[idKey]; exists && len(idValues) > 0 {
			compIngredient.ID = uint(xtremepkg.ToInt(idValues[0]))
			compIngredient.Deleted = xtremepkg.ToBool(formValue[fmt.Sprintf("ingredients[%d][deleted]", i)][0])
			idOK = true
		}

		if childOK || idOK {
			f.Ingredients = append(f.Ingredients, compIngredient)
		} else {
			isLooping = false
		}
	}

	f.Costs = make([]CakeCostForm, 0)
	isLooping = true
	for i := 0; isLooping; i++ {
		var compCost CakeCostForm
		childOK, idOK := false, false

		costTypeKey := fmt.Sprintf("costs[%d][typeId]", i)
		if costTypeValues, exists := formValue[costTypeKey]; exists {
			compCost.TypeId = xtremepkg.ToInt(costTypeValues[0])
			compCost.Price = xtremepkg.ToFloat64(formValue[fmt.Sprintf("costs[%d][price]", i)][0])
			childOK = true
		}

		idKey := fmt.Sprintf("costs[%d][id]", i)
		if idValues, exists := formValue[idKey]; exists && len(idValues) > 0 {
			compCost.ID = uint(xtremepkg.ToInt(idValues[0]))
			compCost.Deleted = xtremepkg.ToBool(formValue[fmt.Sprintf("costs[%d][deleted]", i)][0])
			idOK = true
		}

		if childOK || idOK {
			f.Costs = append(f.Costs, compCost)
		} else {
			isLooping = false
		}
	}
	f.Request = r
}
