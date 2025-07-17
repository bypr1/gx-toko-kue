package handler

import (
	"net/http"

	"service/internal/cake/repository"
	"service/internal/cake/service"
	"service/internal/pkg/constant"
	"service/internal/pkg/core"
	"service/internal/pkg/form"
	cakeparser "service/internal/pkg/parser"

	"github.com/gorilla/mux"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

type CakeHandler struct{}

func (CakeHandler) Get(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewCakeRepository()

	cakes, pagination, _ := repo.Paginate(r.URL.Query())

	psr := cakeparser.CakeParser{Array: cakes}
	res := xtremeres.Response{Array: psr.Briefs(), Pagination: &pagination}
	res.Success(w)
}

func (CakeHandler) Detail(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewCakeRepository()

	cake := repo.FirstById(mux.Vars(r)["id"], repo.PreloadRecipesAndCosts)

	psr := cakeparser.CakeParser{Object: cake}
	res := xtremeres.Response{Object: psr.First(true)}
	res.Success(w)
}

func (CakeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var form form.CakeForm
	form.FormParse(r)
	form.Validate()

	srv := service.NewCakeService()
	cake := srv.Create(form)

	psr := cakeparser.CakeParser{Object: cake}
	res := xtremeres.Response{Object: psr.First(false)}
	res.Success(w)
}

func (CakeHandler) Update(w http.ResponseWriter, r *http.Request) {
	var form form.CakeForm
	form.FormParse(r)
	form.Validate()

	srv := service.NewCakeService()
	cake := srv.Update(form, mux.Vars(r)["id"])

	psr := cakeparser.CakeParser{Object: cake}
	res := xtremeres.Response{Object: psr.First(false)}
	res.Success(w)
}

func (CakeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	srv := service.NewCakeService()
	srv.Delete(mux.Vars(r)["id"])

	res := xtremeres.Response{}
	res.Success(w)
}

func (CakeHandler) GetUnitOfMeasure(w http.ResponseWriter, r *http.Request) {
	result := core.IDName{}.Get(constant.CakeUnitOfMeasure{})

	res := xtremeres.Response{Array: result}
	res.Success(w)
}

func (CakeHandler) GetCostType(w http.ResponseWriter, r *http.Request) {
	result := core.IDName{}.Get(constant.CakeCostType{})

	res := xtremeres.Response{Array: result}
	res.Success(w)
}
