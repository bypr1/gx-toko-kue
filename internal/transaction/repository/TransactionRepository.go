package repository

import (
	"fmt"
	"net/url"
	"service/internal/pkg/config"
	"service/internal/pkg/core"
	errorpkg "service/internal/pkg/error"
	"service/internal/pkg/form"
	"service/internal/pkg/model"
	"service/internal/transaction/excel"

	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	core.TransactionRepository

	core.PaginateRepository[model.Transaction]
	core.FirstIdRepository[model.Transaction]
	core.FindRepository[model.Transaction]
	FindForReport(parameter url.Values) []excel.TransactionReport

	Store(form form.TransactionForm, totalAmount float64) model.Transaction
	Delete(transaction model.Transaction)
	Update(transaction model.Transaction, form form.TransactionForm, totalAmount float64) model.Transaction

	SaveCakes(transaction model.Transaction, form []form.TransactionCakeForm, cakes map[uint]model.Cake) []model.TransactionCake
	DeleteCakes(transaction model.Transaction)

	PreloadCakes(query *gorm.DB) *gorm.DB
}

func NewTransactionRepository(args ...*gorm.DB) TransactionRepository {
	repository := transactionRepository{}
	if len(args) > 0 {
		repository.transaction = args[0]
	}

	return &repository
}

type transactionRepository struct {
	transaction *gorm.DB
}

func (repo *transactionRepository) SetTransaction(tx *gorm.DB) {
	repo.transaction = tx
}

func (repo *transactionRepository) PreloadCakes(query *gorm.DB) *gorm.DB {
	return query.Preload("Cakes.Cake")
}

func (repo *transactionRepository) FirstById(id any, args ...func(query *gorm.DB) *gorm.DB) model.Transaction {
	var transaction model.Transaction
	query := config.PgSQL

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

	query := config.PgSQL.Where("\"createdAt\" BETWEEN ? AND ?", fromDate, toDate)

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

func (repo *transactionRepository) FindForReport(parameter url.Values) []excel.TransactionReport {
	fromDate, toDate := core.SetDateRange(parameter)

	query := config.PgSQL.
		Select(`transactions.*, 
                cakeItems."totalAmount" as "totalAmount", 
                cakeItems."totalCakes" as "totalCakes"`).
		Where("\"createdAt\" BETWEEN ? AND ?", fromDate, toDate).
		Joins(fmt.Sprintf(`
			INNER JOIN (
				SELECT "transactionId", SUM("subTotal") AS "totalAmount", COUNT("id") AS "totalCakes"
				FROM %s
				GROUP BY "transactionId"
			) AS "cakeItems" ON %s."id" = "cakeItems"."transactionId"
		`,
			model.TransactionCake{}.TableName(),
			model.Transaction{}.TableName(),
		)).Preload("Cakes.Cake")

	var transactions []excel.TransactionReport
	err := query.Order("id DESC").Find(&transactions).Error
	if err != nil {
		errorpkg.ErrXtremeTransactionGet(err.Error())
	}

	return transactions
}

func (repo *transactionRepository) Paginate(parameter url.Values) ([]model.Transaction, interface{}, error) {
	fromDate, toDate := core.SetDateRange(parameter)

	query := config.PgSQL.Where("\"createdAt\" BETWEEN ? AND ?", fromDate, toDate)

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

func (repo *transactionRepository) Store(form form.TransactionForm, totalAmount float64) model.Transaction {
	transaction := model.Transaction{
		Date:       core.ParseDate(form.Date),
		TotalPrice: totalAmount,
	}

	err := repo.transaction.Create(&transaction).Error
	if err != nil {
		errorpkg.ErrXtremeTransactionSave(err.Error())
	}

	return transaction
}

func (repo *transactionRepository) Update(transaction model.Transaction, form form.TransactionForm, totalAmount float64) model.Transaction {
	transaction.Date = core.ParseDate(form.Date)
	transaction.TotalPrice = totalAmount

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

func (repo *transactionRepository) SaveCakes(transaction model.Transaction, requests []form.TransactionCakeForm, cakes map[uint]model.Cake) []model.TransactionCake {
	var transactionCakes []model.TransactionCake

	for _, request := range requests {
		var transactionCake model.TransactionCake
		if request.ID > 0 {
			if request.Deleted {
				err := repo.transaction.Where("\"id\" = ?", request.ID).Delete(&model.TransactionCake{}).Error
				if err != nil {
					errorpkg.ErrXtremeTransactionCakeDelete(err.Error())
				}
			} else {
				repo.transaction.Preload("Cake").First(&transactionCake, "id = ?", request.ID)
				if transactionCake.ID == 0 {
					errorpkg.ErrXtremeTransactionCakeGet("ID not found")
				}

				cake := cakes[request.CakeId]
				subTotal := float64(request.Quantity) * cake.Price

				transactionCake.CakeId = request.CakeId
				transactionCake.Quantity = request.Quantity
				transactionCake.Price = cake.Price
				transactionCake.SubTotal = subTotal
				err := repo.transaction.Save(&transactionCake).Error
				if err != nil {
					errorpkg.ErrXtremeTransactionCakeSave(err.Error())
				}
			}
		} else {
			cake := cakes[request.CakeId]
			subTotal := float64(request.Quantity) * cake.Price

			transactionCake = model.TransactionCake{
				TransactionId: transaction.ID,
				CakeId:        request.CakeId,
				Quantity:      request.Quantity,
				Price:         cake.Price,
				SubTotal:      subTotal,
			}
			err := repo.transaction.Save(&transactionCake).Error
			repo.transaction.Preload("Cake").First(&transactionCake, "id = ?", transactionCake.ID)
			if err != nil {
				errorpkg.ErrXtremeTransactionCakeSave(err.Error())
			}
		}

		if transactionCake.ID > 0 {
			transactionCakes = append(transactionCakes, transactionCake)
		}
	}

	return transactionCakes
}

func (repo *transactionRepository) DeleteCakes(transaction model.Transaction) {
	err := repo.transaction.Where("\"transactionId\" = ?", transaction.ID).Delete(&model.TransactionCake{}).Error
	if err != nil {
		errorpkg.ErrXtremeTransactionCakeDelete(err.Error())
	}
}
