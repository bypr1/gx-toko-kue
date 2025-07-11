package handler

import (
	"net/http"

	"service/internal/cake/repository"
	"service/internal/cake/service"
	"service/internal/pkg/form"
	"service/internal/pkg/model"
	cakeparser "service/internal/pkg/parser"

	"github.com/gorilla/mux"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

type IngredientHandler struct{}

func (IngredientHandler) Get(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewCakeComponentIngredientRepository()
	ingredients, pagination, _ := repo.Paginate(r.URL.Query())

	psr := cakeparser.IngredientParser{Array: ingredients}

	res := xtremeres.Response{Array: psr.Get(), Pagination: &pagination}
	res.Success(w)
}

func (IngredientHandler) Detail(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewCakeComponentIngredientRepository()

	var ingredient model.CakeComponentIngredient
	if id := mux.Vars(r)["id"]; id != "" {
		ingredient = repo.FirstById(id)
	}

	psr := cakeparser.IngredientParser{Object: ingredient}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (IngredientHandler) Create(w http.ResponseWriter, r *http.Request) {
	form := form.CakeComponentIngredientForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewIngredientService()
	ingredient := srv.Create(form)

	psr := cakeparser.IngredientParser{Object: ingredient}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (IngredientHandler) Update(w http.ResponseWriter, r *http.Request) {
	form := form.CakeComponentIngredientForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewIngredientService()
	var ingredient model.CakeComponentIngredient
	if id := mux.Vars(r)["id"]; id != "" {
		ingredient = srv.Update(form, id)
	}

	psr := cakeparser.IngredientParser{Object: ingredient}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (IngredientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	srv := service.NewIngredientService()

	if id := mux.Vars(r)["id"]; id != "" {
		srv.Delete(id)
	}

	res := xtremeres.Response{Object: map[string]interface{}{"message": "Ingredient deleted successfully"}}
	res.Success(w)
}
