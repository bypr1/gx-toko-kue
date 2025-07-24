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
	APIPreload(query *gorm.DB) *gorm.DB

	core.PaginateRepository[model.Transaction]
	core.FirstIdRepository[model.Transaction]
	core.FindRepository[model.Transaction]

	FindForReport(form form.TransactionReportForm) []excel.TransactionReport
	Store(form form.TransactionForm, totalAmount float64) model.Transaction
	Delete(transaction model.Transaction)
	Update(transaction model.Transaction, form form.TransactionForm, totalAmount float64) model.Transaction

	SaveCakes(transaction model.Transaction, form []form.TransactionCakeForm, cakes map[uint]model.Cake) []model.TransactionCake
	DeleteCakes(transaction model.Transaction)
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

func (repo *transactionRepository) APIPreload(query *gorm.DB) *gorm.DB {
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

func (repo *transactionRepository) FindForReport(form form.TransactionReportForm) []excel.TransactionReport {
	query := config.PgSQL.
		Select(`transactions.*, 
                cakeItems."totalAmount" as totalAmount, 
                cakeItems."totalCakes" as totalCakes`).
		Where("\"createdAt\" BETWEEN ? AND ?", form.FromDate, form.ToDate).
		Joins(fmt.Sprintf(`
			INNER JOIN (
				SELECT "transactionId", SUM("subTotal") AS "totalAmount", COUNT(id) AS "totalCakes"
				FROM %s
				GROUP BY "transactionId"
			) AS cakeItems ON %s.id = cakeItems."transactionId"
		`,
			model.TransactionCake{}.TableName(),
			model.Transaction{}.TableName(),
		)).Preload("Cakes.Cake")

	if form.MinAmount != nil {
		query = query.Where("\"totalAmount\" >= ?", *form.MinAmount)
	}
	if form.MaxAmount != nil {
		query = query.Where("\"totalAmount\" <= ?", *form.MaxAmount)
	}

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

func (repo *transactionRepository) SaveCakes(transaction model.Transaction, form []form.TransactionCakeForm, cakes map[uint]model.Cake) []model.TransactionCake {
	var transactionCakes []model.TransactionCake

	for _, f := range form {
		var transactionCake model.TransactionCake
		if f.ID > 0 {
			if f.Deleted {
				err := repo.transaction.Where("\"id\" = ?", f.ID).Delete(&model.TransactionCake{}).Error
				if err != nil {
					errorpkg.ErrXtremeTransactionCakeDelete(err.Error())
				}
			} else {
				repo.transaction.Preload("Cake").First(&transactionCake, "id = ?", f.ID)
				if transactionCake.ID == 0 {
					errorpkg.ErrXtremeTransactionCakeGet("ID not found")
				}

				cake := cakes[f.CakeId]
				subTotal := float64(f.Quantity) * cake.Price

				transactionCake.CakeId = f.CakeId
				transactionCake.Quantity = f.Quantity
				transactionCake.Price = cake.Price
				transactionCake.SubTotal = subTotal
				err := repo.transaction.Save(&transactionCake).Error
				if err != nil {
					errorpkg.ErrXtremeTransactionCakeSave(err.Error())
				}
			}
		} else {
			cake := cakes[f.CakeId]
			subTotal := float64(f.Quantity) * cake.Price

			transactionCake = model.TransactionCake{
				TransactionId: transaction.ID,
				CakeId:        f.CakeId,
				Quantity:      f.Quantity,
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
