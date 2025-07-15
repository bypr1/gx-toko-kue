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
	SetCakeRepository(repo port.CakeRepository) *transactionService

	Create(form form.TransactionForm) model.Transaction
	Update(form form.TransactionForm, id string) model.Transaction
	Delete(id string) bool
	DownloadExcel(parameter url.Values) string
}

func NewTransactionService() TransactionService {
	return &transactionService{}
}

type transactionService struct {
	repository     repository.TransactionRepository
	cakeRepository port.CakeRepository
}

func (srv *transactionService) SetCakeRepository(repo port.CakeRepository) *transactionService {
	srv.cakeRepository = repo
	return srv
}

func (srv *transactionService) Create(form form.TransactionForm) model.Transaction {
	var transaction model.Transaction

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.prepareRepository(tx)

		cakes := srv.getCakes(form.Cakes)
		totalAmount := srv.calculateTotalPrice(form.Cakes, cakes)

		transaction = srv.repository.Store(form, totalAmount)
		details := srv.repository.AddCakes(transaction, form.Cakes, cakes)

		transaction.Cakes = append(transaction.Cakes, details...)

		activity.UseActivity{}.SetReference(transaction).SetParser(&parser.TransactionParser{Object: transaction}).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Created new transaction [%d] with total amount %.2f", transaction.ID, transaction.TotalAmount))

		return nil
	})

	return transaction
}

func (srv *transactionService) Update(form form.TransactionForm, id string) model.Transaction {
	var transaction model.Transaction

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.prepareRepository(tx)

		transaction = srv.repository.FirstById(id)

		act := activity.UseActivity{}.
			SetReference(transaction).
			SetParser(&parser.TransactionParser{Object: transaction}).
			SetOldProperty(constant.ACTION_UPDATE)

		cakes := srv.getCakes(form.Cakes)
		totalAmount := srv.calculateTotalPrice(form.Cakes, cakes)

		transaction = srv.repository.Update(transaction, form, totalAmount)
		details := srv.repository.UpdateCakes(transaction, form.Cakes, cakes)

		transaction.Cakes = append(transaction.Cakes, details...)

		act.SetParser(&parser.TransactionParser{Object: transaction}).
			SetNewProperty(constant.ACTION_UPDATE).
			Save(fmt.Sprintf("Updated transaction [%d] with total amount %.2f", transaction.ID, transaction.TotalAmount))

		return nil
	})

	return transaction
}

func (srv *transactionService) Delete(id string) bool {
	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewTransactionRepository(tx)
		transaction := srv.repository.FirstById(id)

		srv.repository.DeleteCakes(transaction)
		srv.repository.Delete(transaction)

		activity.UseActivity{}.SetReference(transaction).SetParser(&parser.TransactionParser{Object: transaction}).SetOldProperty(constant.ACTION_DELETE).
			Save(fmt.Sprintf("Deleted transaction [%d]", transaction.ID))

		return nil
	})
	return true
}

func (srv *transactionService) DownloadExcel(parameter url.Values) string {
	srv.repository = repository.NewTransactionRepository()
	transactions := srv.repository.FindForReport(parameter)

	transactionExcel := excel.TransactionExcel{
		Transactions: transactions,
	}

	filename, _ := transactionExcel.Generate()
	return filename
}

func (srv *transactionService) getCakes(details []form.TransactionCakeForm) map[uint]model.Cake {
	var cakeIDs []any
	for _, detail := range details {
		cakeIDs = append(cakeIDs, detail.CakeID)
	}

	cakes := srv.cakeRepository.FindByIds(cakeIDs)
	cakeMap := make(map[uint]model.Cake)
	for _, cake := range cakes {
		cakeMap[cake.ID] = cake
	}
	return cakeMap
}

func (srv *transactionService) calculateTotalPrice(details []form.TransactionCakeForm, cakeMap map[uint]model.Cake) float64 {
	var totalAmount float64
	for _, detail := range details {
		if cake, exists := cakeMap[detail.CakeID]; exists {
			totalAmount += float64(detail.Quantity) * cake.SellPrice
		}
	}
	return totalAmount
}

func (srv *transactionService) prepareRepository(tx *gorm.DB) {
	if tx != nil {
		srv.repository.SetTransaction(tx)
		srv.cakeRepository.SetTransaction(tx)
	} else {
		srv.repository = repository.NewTransactionRepository(config.PgSQL)
		srv.cakeRepository.SetTransaction(config.PgSQL)
	}
}
