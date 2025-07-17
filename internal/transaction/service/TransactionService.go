package service

import (
	"fmt"
	"net/url"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	"service/internal/pkg/form"
	"service/internal/pkg/model"
	"service/internal/pkg/parser"
	"service/internal/pkg/port"
	"service/internal/transaction/excel"
	"service/internal/transaction/repository"

	"gorm.io/gorm"
)

type TransactionService interface {
	SetCakeRepository(repo port.CakeRepository)

	Create(form form.TransactionForm) model.Transaction
	Update(form form.TransactionForm, id string) model.Transaction
	Delete(id string) bool
	ReportExcel(parameter url.Values) string
}

func NewTransactionService() TransactionService {
	return &transactionService{}
}

type transactionService struct {
	repository     repository.TransactionRepository
	cakeRepository port.CakeRepository
}

func (srv *transactionService) SetCakeRepository(repo port.CakeRepository) {
	srv.cakeRepository = repo
}

func (srv *transactionService) Create(form form.TransactionForm) model.Transaction {
	var transaction model.Transaction
	srv.prepare()

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.setRepositoriesWithTransaction(tx)

		cakes := srv.getCakes(form.Cakes)
		totalAmount := srv.calculateTotalPrice(form.Cakes, cakes)

		transaction = srv.repository.Store(form, totalAmount)
		transaction.Cakes = srv.repository.SaveCakes(transaction, form.Cakes, cakes)

		activity.UseActivity{}.SetReference(transaction).SetParser(&parser.TransactionParser{Object: transaction}).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Created new transaction [%d] with total amount %.2f", transaction.ID, transaction.TotalAmount))

		return nil
	})

	return transaction
}

func (srv *transactionService) Update(form form.TransactionForm, id string) model.Transaction {
	transaction := srv.prepareWithData(id)

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.setRepositoriesWithTransaction(tx)

		act := activity.UseActivity{}.
			SetReference(transaction).
			SetParser(&parser.TransactionParser{Object: transaction}).
			SetOldProperty(constant.ACTION_UPDATE)

		cakes := srv.getCakes(form.Cakes)
		totalAmount := srv.calculateTotalPrice(form.Cakes, cakes)

		transaction = srv.repository.Update(transaction, form, totalAmount)
		transaction.Cakes = srv.repository.SaveCakes(transaction, form.Cakes, cakes)

		act.SetParser(&parser.TransactionParser{Object: transaction}).
			SetNewProperty(constant.ACTION_UPDATE).
			Save(fmt.Sprintf("Updated transaction [%d] with total amount %.2f", transaction.ID, transaction.TotalAmount))

		return nil
	})

	return transaction
}

func (srv *transactionService) Delete(id string) bool {
	transaction := srv.prepareWithData(id)

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.setRepositoriesWithTransaction(tx)

		srv.repository.DeleteCakes(transaction)
		srv.repository.Delete(transaction)

		activity.UseActivity{}.SetReference(transaction).SetParser(&parser.TransactionParser{Object: transaction}).SetOldProperty(constant.ACTION_DELETE).
			Save(fmt.Sprintf("Deleted transaction [%d]", transaction.ID))

		return nil
	})
	return true
}

func (srv *transactionService) ReportExcel(parameter url.Values) string {
	srv.repository = repository.NewTransactionRepository()
	transactions := srv.repository.FindForReport(parameter)

	transactionExcel := excel.TransactionExcel{
		Transactions: transactions,
	}

	filename, _ := transactionExcel.Generate()
	return filename
}

func (srv *transactionService) getCakes(cakeFormItems []form.TransactionCakeForm) map[uint]model.Cake {
	var cakeIDs []any
	for _, cf := range cakeFormItems {
		cakeIDs = append(cakeIDs, cf.CakeID)
	}

	cakes := srv.cakeRepository.FindByIds(cakeIDs)
	cakeMap := make(map[uint]model.Cake)
	for _, cake := range cakes {
		cakeMap[cake.ID] = cake
	}
	return cakeMap
}

func (srv *transactionService) calculateTotalPrice(cakeFormItems []form.TransactionCakeForm, cakeMap map[uint]model.Cake) float64 {
	var totalAmount float64
	for _, cf := range cakeFormItems {
		if cake, exists := cakeMap[cf.CakeID]; exists {
			totalAmount += float64(cf.Quantity) * cake.Price
		}
	}
	return totalAmount
}

func (srv *transactionService) prepare() {
	srv.repository = repository.NewTransactionRepository(nil)
}

func (srv *transactionService) prepareWithData(id any) model.Transaction {
	srv.prepare()
	return srv.repository.FirstById(id)
}

func (srv *transactionService) setRepositoriesWithTransaction(tx *gorm.DB) {
	srv.repository.SetTransaction(tx)
	srv.cakeRepository.SetTransaction(tx)
}
