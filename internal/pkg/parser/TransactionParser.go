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
	transactionObj := parser.Object
	return map[string]interface{}{
		"id":         transactionObj.ID,
		"number":     transactionObj.GetTransactionNumber(),
		"date":       transactionObj.Date.Format("02/01/2006"),
		"totalPrice": transactionObj.TotalPrice,
		"createdAt":  transactionObj.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":  transactionObj.UpdatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser TransactionParser) First() interface{} {
	transactionObj := parser.Object
	var cakes []interface{}
	for _, cake := range transactionObj.Cakes {
		cakes = append(cakes, map[string]interface{}{
			"id":        cake.ID,
			"quantity":  cake.Quantity,
			"price":     cake.Price,
			"subTotal":  cake.SubTotal,
			"cake":      CakeParser{Object: cake.Cake}.Brief(),
			"createdAt": cake.CreatedAt.Format("02/01/2006 15:04"),
			"updatedAt": cake.UpdatedAt.Format("02/01/2006 15:04"),
		})
	}
	return map[string]interface{}{
		"id":         transactionObj.ID,
		"number":     transactionObj.GetTransactionNumber(),
		"date":       transactionObj.Date.Format("02/01/2006"),
		"totalPrice": transactionObj.TotalPrice,
		"createdAt":  transactionObj.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":  transactionObj.UpdatedAt.Format("02/01/2006 15:04"),
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
