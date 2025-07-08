package handler

import (
	"net/http"
	"strconv"

	activityRepo "service/internal/activity/repository"
	"service/internal/cake/repository"
	"service/internal/cake/service"
	form2 "service/internal/pkg/form/cake"
	cakeparser "service/internal/pkg/parser/cake"

	"github.com/gorilla/mux"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

type IngredientHandler struct{}

func (IngredientHandler) Get(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewIngredientRepository()
	ingredients, pagination, _ := repo.Find(r.URL.Query())

	psr := cakeparser.IngredientParser{Array: ingredients}

	res := xtremeres.Response{Array: psr.Get(), Pagination: &pagination}
	res.Success(w)
}

func (IngredientHandler) Detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	repo := repository.NewIngredientRepository()
	ingredient := repo.FirstById(uint(id))

	psr := cakeparser.IngredientParser{Object: ingredient}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (IngredientHandler) Create(w http.ResponseWriter, r *http.Request) {
	form := form2.IngredientForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewIngredientService()
	srv.SetActivityRepository(activityRepo.NewActivityRepository())

	ingredient := srv.Create(form)

	psr := cakeparser.IngredientParser{Object: ingredient}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (IngredientHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	form := form2.IngredientForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewIngredientService()
	srv.SetActivityRepository(activityRepo.NewActivityRepository())

	ingredient := srv.Update(form, uint(id))

	psr := cakeparser.IngredientParser{Object: ingredient}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (IngredientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	srv := service.NewIngredientService()
	srv.SetActivityRepository(activityRepo.NewActivityRepository())

	err := srv.Delete(uint(id))
	if err != nil {
		xtremeres.Error(http.StatusInternalServerError, "Unable to delete ingredient", err.Error(), false, nil)
		return
	}

	res := xtremeres.Response{Object: map[string]interface{}{"message": "Ingredient deleted successfully"}}
	res.Success(w)
}
