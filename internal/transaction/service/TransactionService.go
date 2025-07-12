package service

import (
	"fmt"
	"service/internal/cake/repository"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	"service/internal/pkg/form"
	"service/internal/pkg/model"
	"service/internal/pkg/parser"
	"service/internal/pkg/port"
	transactionRepository "service/internal/transaction/repository"

	"gorm.io/gorm"
)

type TransactionService interface {
	Create(form form.TransactionForm) model.Transaction
	Update(form form.TransactionForm, id string) model.Transaction
	Delete(id string) bool
}

func NewTransactionService() TransactionService {
	return &transactionService{}
}

type transactionService struct {
	repository     transactionRepository.TransactionRepository
	cakeRepository port.CakeRepository
}

func (srv *transactionService) Create(form form.TransactionForm) model.Transaction {
	var transaction model.Transaction

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = transactionRepository.NewTransactionRepository(tx)
		srv.cakeRepository = repository.NewCakeRepository(tx)

		cakes, totalAmount := srv.getCakesAndCalculateTotal(form.Details)

		transaction = srv.repository.Store(form, totalAmount)
		details := srv.repository.AddDetails(transaction, form.Details, cakes)

		transaction.Details = append(transaction.Details, details...)

		activity.UseActivity{}.SetReference(transaction).SetParser(&parser.TransactionParser{Object: transaction}).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Created new transaction [%d] with total amount %.2f", transaction.ID, transaction.TotalAmount))

		return nil
	})

	return transaction
}

func (srv *transactionService) Update(form form.TransactionForm, id string) model.Transaction {
	var transaction model.Transaction

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = transactionRepository.NewTransactionRepository(tx)
		srv.cakeRepository = repository.NewCakeRepository(tx)
		transaction = srv.repository.FirstById(id)

		act := activity.UseActivity{}.
			SetReference(transaction).
			SetParser(&parser.TransactionParser{Object: transaction}).
			SetOldProperty(constant.ACTION_UPDATE)

		cakes, totalAmount := srv.getCakesAndCalculateTotal(form.Details)

		transaction = srv.repository.Update(transaction, form, totalAmount)
		details := srv.repository.UpdateDetails(transaction, form.Details, cakes)

		transaction.Details = append(transaction.Details, details...)

		act.SetParser(&parser.TransactionParser{Object: transaction}).
			SetNewProperty(constant.ACTION_UPDATE).
			Save(fmt.Sprintf("Updated transaction [%d] with total amount %.2f", transaction.ID, transaction.TotalAmount))

		return nil
	})

	return transaction
}

func (srv *transactionService) Delete(id string) bool {
	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = transactionRepository.NewTransactionRepository(tx)
		transaction := srv.repository.FirstById(id)

		srv.repository.DeleteDetails(transaction)
		srv.repository.Delete(transaction)

		activity.UseActivity{}.SetReference(transaction).SetParser(&parser.TransactionParser{Object: transaction}).SetOldProperty(constant.ACTION_DELETE).
			Save(fmt.Sprintf("Deleted transaction [%d]", transaction.ID))

		return nil
	})
	return true
}

func (srv *transactionService) getCakesAndCalculateTotal(details []form.TransactionDetailCakeForm) (map[uint]model.Cake, float64) {
	var cakeIDs []any
	for _, detail := range details {
		cakeIDs = append(cakeIDs, detail.CakeID)
	}

	cakes := srv.cakeRepository.FindByIds(cakeIDs)
	cakeMap := make(map[uint]model.Cake)
	var totalAmount float64

	for _, cake := range cakes {
		cakeMap[cake.ID] = cake
	}

	for _, detail := range details {
		if cake, exists := cakeMap[detail.CakeID]; exists {
			totalAmount += float64(detail.Quantity) * cake.SellPrice
		}
	}

	return cakeMap, totalAmount
}
