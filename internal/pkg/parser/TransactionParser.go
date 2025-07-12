package parser

import (
	"service/internal/pkg/model"
)

type TransactionParser struct {
	Array  []model.Transaction
	Object model.Transaction
}

func (parser TransactionParser) CreateActivity(action string) interface{} {
	return parser.Brief()
}

func (parser TransactionParser) DeleteActivity(action string) interface{} {
	return parser.Brief()
}

func (parser TransactionParser) GeneralActivity(action string) interface{} {
	return parser.Brief()
}

func (parser TransactionParser) UpdateActivity(action string) interface{} {
	return parser.Brief()
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
	var details []interface{}
	for _, detail := range transactionObj.Details {
		details = append(details, TransactionDetailCakeParser{Object: detail}.Brief())
	}
	return map[string]interface{}{
		"id":              transactionObj.ID,
		"transactionDate": transactionObj.TransactionDate.Format("02/01/2006"),
		"totalAmount":     transactionObj.TotalAmount,
		"createdAt":       transactionObj.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":       transactionObj.UpdatedAt.Format("02/01/2006 15:04"),
		"details":         details,
	}
}

func (parser TransactionParser) Get() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, TransactionParser{Object: obj}.First())
	}
	return result
}

func (parser TransactionParser) First() interface{} {
	transactionObj := parser.Object
	return map[string]interface{}{
		"id":              transactionObj.ID,
		"transactionDate": transactionObj.TransactionDate.Format("02/01/2006"),
		"totalAmount":     transactionObj.TotalAmount,
		"createdAt":       transactionObj.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":       transactionObj.UpdatedAt.Format("02/01/2006 15:04"),
	}
}
