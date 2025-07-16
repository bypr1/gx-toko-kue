package handler

import (
	"net/http"
	cakeRepository "service/internal/cake/repository"
	"service/internal/pkg/form"
	"service/internal/pkg/parser"
	transactionRepository "service/internal/transaction/repository"
	transactionService "service/internal/transaction/service"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"github.com/gorilla/mux"
)

type TransactionHandler struct{}

func (TransactionHandler) Get(w http.ResponseWriter, r *http.Request) {
	repo := transactionRepository.NewTransactionRepository()
	transactions, pagination, _ := repo.Paginate(r.URL.Query())

	transactionParser := parser.TransactionParser{Array: transactions}
	res := xtremeres.Response{Array: transactionParser.Briefs(), Pagination: &pagination}
	res.Success(w)
}

func (TransactionHandler) Detail(w http.ResponseWriter, r *http.Request) {
	repo := transactionRepository.NewTransactionRepository()
	transaction := repo.FirstById(mux.Vars(r)["id"], repo.PreloadCakes)

	transactionParser := parser.TransactionParser{Object: transaction}
	res := xtremeres.Response{Object: transactionParser.First()}
	res.Success(w)
}

func (TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var transactionForm form.TransactionForm
	transactionForm.APIParse(r)
	transactionForm.Validate()

	service := transactionService.NewTransactionService()
	service.SetCakeRepository(cakeRepository.NewCakeRepository())
	transaction := service.Create(transactionForm)

	transactionParser := parser.TransactionParser{Object: transaction}
	res := xtremeres.Response{Object: transactionParser.First()}
	res.Success(w)
}

func (TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	var transactionForm form.TransactionForm
	transactionForm.APIParse(r)
	transactionForm.Validate()

	service := transactionService.NewTransactionService()
	service.SetCakeRepository(cakeRepository.NewCakeRepository())
	transaction := service.Update(transactionForm, mux.Vars(r)["id"])

	transactionParser := parser.TransactionParser{Object: transaction}
	res := xtremeres.Response{Object: transactionParser.First()}
	res.Success(w)
}

func (TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	service := transactionService.NewTransactionService()
	service.Delete(mux.Vars(r)["id"])

	res := xtremeres.Response{}
	res.Success(w)
}

func (TransactionHandler) ReportExcel(w http.ResponseWriter, r *http.Request) {
	service := transactionService.NewTransactionService()
	filename := service.ReportExcel(r.URL.Query())

	res := xtremeres.Response{Object: filename}
	res.Success(w)
}
