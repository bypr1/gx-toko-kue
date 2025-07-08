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
	"gorm.io/gorm"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

type CakeHandler struct{}

func (CakeHandler) Get(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewCakeRepository()
	cakes, pagination, _ := repo.Find(r.URL.Query())

	psr := cakeparser.CakeParser{Array: cakes}

	res := xtremeres.Response{Array: psr.Get(), Pagination: &pagination}
	res.Success(w)
}

func (CakeHandler) Detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	repo := repository.NewCakeRepository()
	cake := repo.FirstById(id)

	psr := cakeparser.CakeParser{Object: cake}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (CakeHandler) Create(w http.ResponseWriter, r *http.Request) {
	form := form2.CakeForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewCakeService()
	srv.SetActivityRepository(activityRepo.NewActivityRepository())

	cake := srv.Create(form)

	psr := cakeparser.CakeParser{Object: cake}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (CakeHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	form := form2.CakeForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewCakeService()
	srv.SetActivityRepository(activityRepo.NewActivityRepository())

	cake := srv.Update(form, uint(id))

	psr := cakeparser.CakeParser{Object: cake}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (CakeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	srv := service.NewCakeService()
	srv.SetActivityRepository(activityRepo.NewActivityRepository())

	err := srv.Delete(uint(id))
	if err != nil {
		xtremeres.Error(http.StatusInternalServerError, "Unable to delete cake", err.Error(), false, nil)
		return
	}

	res := xtremeres.Response{Object: map[string]interface{}{"message": "Cake deleted successfully"}}
	res.Success(w)
}

func (CakeHandler) CalculateCost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	repo := repository.NewCakeRepository()
	cake := repo.FirstById(uint(id), func(query *gorm.DB) *gorm.DB {
		return query.Preload("Recipes").Preload("Recipes.Ingredient").Preload("Costs")
	})

	srv := service.NewCakeService()
	totalCost := srv.CalculateCakeCost(cake)

	psr := cakeparser.CakeCostCalculationParser{
		Cake:      cake,
		TotalCost: totalCost,
	}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}
