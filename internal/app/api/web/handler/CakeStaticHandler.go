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
	var uom constant.UnitOfMeasure

	if isIngredient := xtremepkg.ToBool(r.URL.Query().Get("isIngredient")); isIngredient {
		uom = uom.SetIsIngredient()
	} else if isCake := xtremepkg.ToBool(r.URL.Query().Get("isCake")); isCake {
		uom = uom.SetIsCake()
	}
	result = core.IDName{}.Get(uom)

	res := xtremeres.Response{Array: result}
	res.Success(w)
}

func (CakeStaticHandler) GetCostType(w http.ResponseWriter, r *http.Request) {
	result := core.IDName{}.Get(constant.CakeCostType{})

	res := xtremeres.Response{Array: result}
	res.Success(w)
}
