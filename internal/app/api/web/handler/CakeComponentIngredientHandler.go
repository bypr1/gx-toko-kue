package handler

import (
	"net/http"

	"service/internal/cake/repository"
	"service/internal/cake/service"
	"service/internal/pkg/form"
	cakeparser "service/internal/pkg/parser"

	"github.com/gorilla/mux"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

type IngredientHandler struct{}

func (IngredientHandler) Get(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewCakeComponentIngredientRepository()
	ingredients, pagination, _ := repo.Paginate(r.URL.Query())

	psr := cakeparser.CakeComponentIngredientParser{Array: ingredients}

	res := xtremeres.Response{Array: psr.Get(), Pagination: &pagination}
	res.Success(w)
}

func (IngredientHandler) Detail(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewCakeComponentIngredientRepository()

	ingredient := repo.FirstById(mux.Vars(r)["id"])

	psr := cakeparser.CakeComponentIngredientParser{Object: ingredient}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (IngredientHandler) Create(w http.ResponseWriter, r *http.Request) {
	form := form.CakeComponentIngredientForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewIngredientService()
	ingredient := srv.Create(form)

	psr := cakeparser.CakeComponentIngredientParser{Object: ingredient}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (IngredientHandler) Update(w http.ResponseWriter, r *http.Request) {
	form := form.CakeComponentIngredientForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewIngredientService()
	ingredient := srv.Update(form, mux.Vars(r)["id"])

	psr := cakeparser.CakeComponentIngredientParser{Object: ingredient}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (IngredientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	srv := service.NewIngredientService()
	srv.Delete(mux.Vars(r)["id"])

	res := xtremeres.Response{}
	res.Success(w)
}
