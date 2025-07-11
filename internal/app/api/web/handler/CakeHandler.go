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

	var cake model.Cake
	if id := mux.Vars(r)["id"]; id != "" {
		cake = repo.FirstById(id)
	}

	psr := cakeparser.CakeParser{Object: cake}
	res := xtremeres.Response{Object: psr.Brief()}
	res.Success(w)
}

func (CakeHandler) Create(w http.ResponseWriter, r *http.Request) {
	form := form.CakeForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewCakeService()
	cake := srv.Create(form)

	psr := cakeparser.CakeParser{Object: cake}
	res := xtremeres.Response{Object: psr.Brief()}
	res.Success(w)
}

func (CakeHandler) Update(w http.ResponseWriter, r *http.Request) {
	form := form.CakeForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewCakeService()

	var cake model.Cake
	if id := mux.Vars(r)["id"]; id != "" {
		cake = srv.Update(form, id)
	}

	psr := cakeparser.CakeParser{Object: cake}
	res := xtremeres.Response{Object: psr.Brief()}
	res.Success(w)
}

func (CakeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	srv := service.NewCakeService()

	if id := mux.Vars(r)["id"]; id != "" {
		srv.Delete(id)
	}

	res := xtremeres.Response{Object: map[string]interface{}{"message": "Cake deleted successfully"}}
	res.Success(w)
}
