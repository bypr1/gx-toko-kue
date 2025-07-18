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

// -- Preload necessary relations for the model ---

func (repo *transactionRepository) PreloadCakes(query *gorm.DB) *gorm.DB {
	return query.Preload("Cakes.Cake")
}

// -- Public operations that interact with the database ---

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
		TransactionDate: core.ParseDate(form.TransactionDate),
		TotalAmount:     totalAmount,
	}

	err := repo.transaction.Create(&transaction).Error
	if err != nil {
		errorpkg.ErrXtremeTransactionSave(err.Error())
	}

	return transaction
}

func (repo *transactionRepository) Update(transaction model.Transaction, form form.TransactionForm, totalAmount float64) model.Transaction {
	transaction.TransactionDate = core.ParseDate(form.TransactionDate)
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

func (repo *transactionRepository) DeleteCakes(transaction model.Transaction) {
	err := repo.transaction.Where("\"transactionId\" = ?", transaction.ID).Delete(&model.TransactionCake{}).Error
	if err != nil {
		errorpkg.ErrXtremeTransactionCakeDelete(err.Error())
	}
}

func (repo *transactionRepository) SaveCakes(transaction model.Transaction, form []form.TransactionCakeForm, cakes map[uint]model.Cake) []model.TransactionCake {
	var toDelete []uint
	var toUpdate []model.TransactionCake
	var toCreate []model.TransactionCake

	existingCakes := repo.getExistingCakes(transaction)
	existingMap := make(map[uint]model.TransactionCake)
	for _, existing := range existingCakes {
		existingMap[existing.ID] = existing
	}

	for _, f := range form {
		if f.Deleted {
			toDelete = append(toDelete, f.ID)
		} else {
			cake := cakes[f.CakeID]
			subTotal := float64(f.Quantity) * cake.Price

			transactionCake := model.TransactionCake{
				TransactionId: transaction.ID,
				CakeId:        f.CakeID,
				Quantity:      f.Quantity,
				Price:         cake.Price,
				SubTotal:      subTotal,
			}

			if f.ID > 0 { // Update existing
				transactionCake.ID = f.ID
				toUpdate = append(toUpdate, transactionCake)
			} else { // Create new
				toCreate = append(toCreate, transactionCake)
			}
		}
	}

	if err := repo.batchDeleteCakes(toDelete); err != nil {
		errorpkg.ErrXtremeTransactionCakeSave(err.Error())
	}

	if err := repo.batchUpdateCakes(toUpdate); err != nil {
		errorpkg.ErrXtremeTransactionCakeSave(err.Error())
	}

	if err := repo.batchCreateCakes(toCreate); err != nil {
		errorpkg.ErrXtremeTransactionCakeSave(err.Error())
	}

	return append(toUpdate, toCreate...)
}

// -- Private helper sections for the repository ---

func (repo *transactionRepository) getExistingCakes(transaction model.Transaction) []model.TransactionCake {
	var cakes []model.TransactionCake
	repo.transaction.Where("\"transactionId\" = ?", transaction.ID).Find(&cakes)
	return cakes
}

func (repo *transactionRepository) batchDeleteCakes(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	return repo.transaction.Where("\"id\" IN ?", ids).Delete(&model.TransactionCake{}).Error
}

func (repo *transactionRepository) batchUpdateCakes(cakes []model.TransactionCake) error {
	if len(cakes) == 0 {
		return nil
	}

	for i := range cakes {
		repo.transaction.Save(&cakes[i])
		repo.transaction.Preload("Cake").First(&cakes[i], cakes[i].ID)
	}
	return nil
}

func (repo *transactionRepository) batchCreateCakes(cakes []model.TransactionCake) error {
	if len(cakes) == 0 {
		return nil
	}

	repo.transaction.Create(&cakes)
	for i := range cakes {
		repo.transaction.Preload("Cake").First(&cakes[i], cakes[i].ID)
	}

	return nil
}
