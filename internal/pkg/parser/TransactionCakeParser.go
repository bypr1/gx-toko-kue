package parser

import (
	"service/internal/pkg/model"
)

type TransactionCakeParser struct {
	Array  []model.TransactionCake
	Object model.TransactionCake
}

func (parser TransactionCakeParser) CreateActivity(action string) interface{} {
	return parser.Brief()
}

func (parser TransactionCakeParser) DeleteActivity(action string) interface{} {
	return parser.Brief()
}

func (parser TransactionCakeParser) GeneralActivity(action string) interface{} {
	return parser.Brief()
}

func (parser TransactionCakeParser) UpdateActivity(action string) interface{} {
	return parser.Brief()
}

func (parser TransactionCakeParser) Get() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, TransactionCakeParser{Object: obj}.First())
	}
	return result
}

func (parser TransactionCakeParser) First() interface{} {
	txCakeObj := parser.Object
	return map[string]interface{}{
		"id":            txCakeObj.ID,
		"transactionId": txCakeObj.TransactionID,
		"quantity":      txCakeObj.Quantity,
		"price":         txCakeObj.Price,
		"subTotal":      txCakeObj.SubTotal,
		"cake":          CakeParser{Object: txCakeObj.Cake}.Brief(),
		"createdAt":     txCakeObj.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":     txCakeObj.UpdatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser TransactionCakeParser) Briefs() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, TransactionCakeParser{Object: obj}.Brief())
	}
	return result
}

func (parser TransactionCakeParser) Brief() interface{} {
	txCakeObj := parser.Object
	return map[string]interface{}{
		"id":        txCakeObj.ID,
		"cakeId":    txCakeObj.CakeID,
		"quantity":  txCakeObj.Quantity,
		"price":     txCakeObj.Price,
		"subTotal":  txCakeObj.SubTotal,
		"createdAt": txCakeObj.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt": txCakeObj.UpdatedAt.Format("02/01/2006 15:04"),
	}
}
