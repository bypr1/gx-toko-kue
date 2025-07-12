package handler

import (
	"net/http"
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
	vars := mux.Vars(r)
	id := vars["id"]

	repo := transactionRepository.NewTransactionRepository()
	transaction := repo.FirstById(id, repo.WithDetails)

	transactionParser := parser.TransactionParser{Object: transaction}
	res := xtremeres.Response{Object: transactionParser.Brief()}
	res.Success(w)
}

func (TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var transactionForm form.TransactionForm
	transactionForm.APIParse(r)
	transactionForm.Validate()

	service := transactionService.NewTransactionService()
	transaction := service.Create(transactionForm)

	transactionParser := parser.TransactionParser{Object: transaction}
	res := xtremeres.Response{Object: transactionParser.Brief()}
	res.Success(w)
}

func (TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var transactionForm form.TransactionForm
	transactionForm.APIParse(r)
	transactionForm.Validate()

	service := transactionService.NewTransactionService()
	transaction := service.Update(transactionForm, id)

	transactionParser := parser.TransactionParser{Object: transaction}
	res := xtremeres.Response{Object: transactionParser.Brief()}
	res.Success(w)
}

func (TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	service := transactionService.NewTransactionService()
	service.Delete(id)

	res := xtremeres.Response{}
	res.Success(w)
}
