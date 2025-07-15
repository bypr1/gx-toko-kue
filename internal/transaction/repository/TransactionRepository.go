package repository

import (
	"fmt"
	"net/url"
	"service/internal/pkg/config"
	"service/internal/pkg/core"
	errorpkg "service/internal/pkg/error"
	formpkg "service/internal/pkg/form"
	"service/internal/pkg/model"

	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	core.TransactionRepository

	core.PaginateRepository[model.Transaction]
	core.FirstIdRepository[model.Transaction]
	core.FindRepository[model.Transaction]

	FindForReport(parameter url.Values) []model.TransactionReport

	Store(form formpkg.TransactionForm, totalAmount float64) model.Transaction
	Delete(transaction model.Transaction)
	Update(transaction model.Transaction, form formpkg.TransactionForm, totalAmount float64) model.Transaction

	AddCakes(transaction model.Transaction, details []formpkg.TransactionDetailCakeForm, cakes map[uint]model.Cake) []model.TransactionDetailCake
	UpdateCakes(transaction model.Transaction, details []formpkg.TransactionDetailCakeForm, cakes map[uint]model.Cake) []model.TransactionDetailCake
	DeleteCakes(transaction model.Transaction)

	PreloadCakes(query *gorm.DB) *gorm.DB
}

func NewTransactionRepository(args ...*gorm.DB) TransactionRepository {
	repository := transactionRepository{}
	if len(args) > 0 {
		repository.transaction = args[0]
	} else {
		repository.transaction = config.PgSQL // Default to global config
	}

	return &repository
}

type transactionRepository struct {
	transaction *gorm.DB
}

func (repo *transactionRepository) SetTransaction(tx *gorm.DB) {
	repo.transaction = tx
}

func (repo *transactionRepository) FirstById(id any, args ...func(query *gorm.DB) *gorm.DB) model.Transaction {
	var transaction model.Transaction
	query := repo.transaction

	if len(args) > 0 {
		query = args[0](query)
	}

	err := query.First(&transaction, "id = ?", id).Error
	if err != nil {
		errorpkg.ErrXtremeTransactionGet(err.Error())
	}

	return transaction
}

func (repo *transactionRepository) Find(parameter url.Values) []model.Transaction {
	fromDate, toDate := core.SetDateRange(parameter)

	query := repo.transaction.Where("\"createdAt\" BETWEEN ? AND ?", fromDate, toDate)

	if minAmount := parameter.Get("minAmount"); minAmount != "" {
		query = query.Where("\"totalAmount\" >= ?", minAmount)
	}

	if maxAmount := parameter.Get("maxAmount"); maxAmount != "" {
		query = query.Where("\"totalAmount\" <= ?", maxAmount)
	}

	var transactions []model.Transaction
	err := query.Order("id DESC").Find(&transactions).Error
	if err != nil {
		errorpkg.ErrXtremeTransactionGet(err.Error())
	}

	return transactions
}

func (repo *transactionRepository) FindForReport(parameter url.Values) []model.TransactionReport {
	fromDate, toDate := core.SetDateRange(parameter)

	query := repo.transaction.
		Select(`transactions.*, 
                details."totalAmount" as "totalAmount", 
                details."totalCakes" as "totalCakes"`).
		Where("\"createdAt\" BETWEEN ? AND ?", fromDate, toDate).
		Joins(fmt.Sprintf(`
			INNER JOIN (
				SELECT "transactionId", SUM("subTotal") AS "totalAmount", COUNT("id") AS "totalCakes"
				FROM %s
				GROUP BY "transactionId"
			) AS "details" ON %s."id" = "details"."transactionId"
		`,
			model.TransactionDetailCake{}.TableName(),
			model.Transaction{}.TableName(),
		)).Preload("Details.Cake")

	var transactions []model.TransactionReport
	err := query.Order("id DESC").Find(&transactions).Error
	if err != nil {
		errorpkg.ErrXtremeTransactionGet(err.Error())
	}

	return transactions
}

func (repo *transactionRepository) Paginate(parameter url.Values) ([]model.Transaction, interface{}, error) {
	fromDate, toDate := core.SetDateRange(parameter)

	query := repo.transaction.Where("\"createdAt\" BETWEEN ? AND ?", fromDate, toDate)

	if minAmount := parameter.Get("minAmount"); minAmount != "" {
		query = query.Where("\"totalAmount\" >= ?", minAmount)
	}

	if maxAmount := parameter.Get("maxAmount"); maxAmount != "" {
		query = query.Where("\"totalAmount\" <= ?", maxAmount)
	}

	transactions, pagination, err := xtrememodel.Paginate(query.Order("id DESC"), parameter, model.Transaction{})
	if err != nil {
		return nil, nil, err
	}

	return transactions, pagination, nil
}

func (repo *transactionRepository) Store(form formpkg.TransactionForm, totalAmount float64) model.Transaction {
	transaction := model.Transaction{
		TransactionDate: form.GetTransactionDate(),
		TotalAmount:     totalAmount,
	}

	err := repo.transaction.Create(&transaction).Error
	if err != nil {
		errorpkg.ErrXtremeTransactionSave(err.Error())
	}

	return transaction
}

func (repo *transactionRepository) Update(transaction model.Transaction, form formpkg.TransactionForm, totalAmount float64) model.Transaction {
	transaction.TransactionDate = form.GetTransactionDate()
	transaction.TotalAmount = totalAmount

	err := repo.transaction.Save(&transaction).Error
	if err != nil {
		errorpkg.ErrXtremeTransactionSave(err.Error())
	}

	return transaction
}

func (repo *transactionRepository) Delete(transaction model.Transaction) {
	err := repo.transaction.Delete(&transaction).Error
	if err != nil {
		errorpkg.ErrXtremeTransactionDelete(err.Error())
	}
}

func (repo *transactionRepository) addDetail(transaction model.Transaction, detail formpkg.TransactionDetailCakeForm, cake model.Cake) model.TransactionDetailCake {
	subTotal := float64(detail.Quantity) * cake.SellPrice

	transactionDetail := model.TransactionDetailCake{
		TransactionID: transaction.ID,
		CakeID:        detail.CakeID,
		Quantity:      detail.Quantity,
		UnitPrice:     cake.SellPrice,
		SubTotal:      subTotal,
	}

	err := repo.transaction.Create(&transactionDetail).Error
	if err != nil {
		errorpkg.ErrXtremeTransactionDetailSave(err.Error())
	}

	return transactionDetail
}

func (repo *transactionRepository) AddCakes(transaction model.Transaction, details []formpkg.TransactionDetailCakeForm, cakes map[uint]model.Cake) []model.TransactionDetailCake {
	var transactionDetails []model.TransactionDetailCake
	for _, detail := range details {
		cake := cakes[detail.CakeID]
		transactionDetail := repo.addDetail(transaction, detail, cake)
		transactionDetails = append(transactionDetails, transactionDetail)
	}
	return transactionDetails
}

func (repo *transactionRepository) UpdateCakes(transaction model.Transaction, details []formpkg.TransactionDetailCakeForm, cakes map[uint]model.Cake) []model.TransactionDetailCake {
	repo.DeleteCakes(transaction)
	return repo.AddCakes(transaction, details, cakes)
}

func (repo *transactionRepository) DeleteCakes(transaction model.Transaction) {
	err := repo.transaction.Where("\"transactionId\" = ?", transaction.ID).Delete(&model.TransactionDetailCake{}).Error
	if err != nil {
		errorpkg.ErrXtremeTransactionDetailDelete(err.Error())
	}
}

func (repo *transactionRepository) PreloadCakes(query *gorm.DB) *gorm.DB {
	return query.Preload("Cakes.Cake")
}
