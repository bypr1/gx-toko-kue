package parser

import (
	"service/internal/pkg/model"
)

type TransactionParser struct {
	Array  []model.Transaction
	Object model.Transaction
}

func (parser TransactionParser) Briefs() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, TransactionParser{Object: obj}.Brief())
	}
	return result
}

func (parser TransactionParser) Brief() interface{} {
	transaction := parser.Object
	return map[string]interface{}{
		"id":         transaction.ID,
		"number":     transaction.GetTransactionNumber(),
		"date":       transaction.Date.Format("02/01/2006"),
		"totalPrice": transaction.TotalPrice,
		"createdAt":  transaction.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":  transaction.UpdatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser TransactionParser) First() interface{} {
	transaction := parser.Object
	var cakes []interface{}
	for _, cake := range transaction.Cakes {
		cakes = append(cakes, map[string]interface{}{
			"id":       cake.ID,
			"quantity": cake.Quantity,
			"price":    cake.Price,
			"subTotal": cake.SubTotal,
			"cake":     CakeParser{Object: cake.Cake}.Brief(),
		})
	}
	return map[string]interface{}{
		"id":         transaction.ID,
		"number":     transaction.GetTransactionNumber(),
		"date":       transaction.Date.Format("02/01/2006"),
		"totalPrice": transaction.TotalPrice,
		"createdAt":  transaction.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":  transaction.UpdatedAt.Format("02/01/2006 15:04"),
		"cakes":      cakes,
	}
}

func (parser TransactionParser) CreateActivity(action string) interface{} {
	return parser.Brief()
}

func (parser TransactionParser) UpdateActivity(action string) interface{} {
	return parser.Brief()
}

func (parser TransactionParser) DeleteActivity(action string) interface{} {
	return parser.Brief()
}

func (parser TransactionParser) GeneralActivity(action string) interface{} {
	return parser.Brief()
}
