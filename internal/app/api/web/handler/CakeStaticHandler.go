package handler

import (
	"net/http"

	"service/internal/pkg/constant"
	"service/internal/pkg/core"

	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

type CakeStaticHandler struct{}

func (CakeStaticHandler) GetUnitOfMeasure(w http.ResponseWriter, r *http.Request) {
	var result []interface{}

	if xtremepkg.ToBool(r.URL.Query().Get("isIngredient")) {
		result = core.IDName{}.Get(constant.CakeIngredientUnitOfMeasure{})
	} else if xtremepkg.ToBool(r.URL.Query().Get("isCake")) {
		result = core.IDName{}.Get(constant.CakeUnitOfMeasure{})
	} else {
		result = core.IDName{}.Get(constant.UnitOfMeasure{})
	}

	res := xtremeres.Response{Array: result}
	res.Success(w)
}

func (CakeStaticHandler) GetCostType(w http.ResponseWriter, r *http.Request) {
	result := core.IDName{}.Get(constant.CakeCostType{})

	res := xtremeres.Response{Array: result}
	res.Success(w)
}
