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
		cakes = append(cakes, TransactionCakeParser{Object: cake}.First())
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

type TransactionCakeParser struct {
	Array  []model.TransactionCake
	Object model.TransactionCake
}

func (parser TransactionCakeParser) Brief() interface{} {
	txCakeObj := parser.Object
	return map[string]interface{}{
		"id":        txCakeObj.ID,
		"cakeId":    txCakeObj.CakeId,
		"quantity":  txCakeObj.Quantity,
		"price":     txCakeObj.Price,
		"subTotal":  txCakeObj.SubTotal,
		"createdAt": txCakeObj.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt": txCakeObj.UpdatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser TransactionCakeParser) First() interface{} {
	txCakeObj := parser.Object
	return map[string]interface{}{
		"id":            txCakeObj.ID,
		"transactionId": txCakeObj.TransactionId,
		"quantity":      txCakeObj.Quantity,
		"price":         txCakeObj.Price,
		"subTotal":      txCakeObj.SubTotal,
		"cake":          CakeParser{Object: txCakeObj.Cake}.Brief(),
		"createdAt":     txCakeObj.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":     txCakeObj.UpdatedAt.Format("02/01/2006 15:04"),
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
