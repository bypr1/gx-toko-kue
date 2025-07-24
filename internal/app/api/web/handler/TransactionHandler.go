package handler

import (
	"net/http"
	cakeRepository "service/internal/cake/repository"
	"service/internal/pkg/form"
	"service/internal/pkg/parser"
	"service/internal/transaction/repository"
	"service/internal/transaction/service"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"github.com/gorilla/mux"
)

type TransactionHandler struct{}

func (TransactionHandler) Get(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewTransactionRepository()
	transactions, pagination, _ := repo.Paginate(r.URL.Query())

	psr := parser.TransactionParser{Array: transactions}
	res := xtremeres.Response{Array: psr.Briefs(), Pagination: &pagination}
	res.Success(w)
}

func (TransactionHandler) Detail(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewTransactionRepository()
	transaction := repo.FirstById(mux.Vars(r)["id"], repo.APIPreload)

	psr := parser.TransactionParser{Object: transaction}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var form form.TransactionForm
	form.APIParse(r)
	form.Validate()

	srv := service.NewTransactionService()
	srv.SetCakeRepository(cakeRepository.NewCakeRepository())
	transaction := srv.Create(form)

	psr := parser.TransactionParser{Object: transaction}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	var form form.TransactionForm
	form.APIParse(r)
	form.Validate()

	srv := service.NewTransactionService()
	srv.SetCakeRepository(cakeRepository.NewCakeRepository())
	transaction := srv.Update(form, mux.Vars(r)["id"])

	psr := parser.TransactionParser{Object: transaction}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	srv := service.NewTransactionService()
	srv.SetCakeRepository(cakeRepository.NewCakeRepository())
	srv.Delete(mux.Vars(r)["id"])

	res := xtremeres.Response{}
	res.Success(w)
}

func (TransactionHandler) ReportExcel(w http.ResponseWriter, r *http.Request) {
	var form form.TransactionReportForm
	form.APIParse(r)
	form.Validate()

	srv := service.NewTransactionService()
	srv.SetCakeRepository(cakeRepository.NewCakeRepository())
	filename := srv.ReportExcel(form)

	res := xtremeres.Response{Object: filename}
	res.Success(w)
}
