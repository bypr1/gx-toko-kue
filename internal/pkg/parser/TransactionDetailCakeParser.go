package parser

import (
	"service/internal/pkg/model"
)

type TransactionDetailCakeParser struct {
	Array  []model.TransactionDetailCake
	Object model.TransactionDetailCake
}

func (parser TransactionDetailCakeParser) CreateActivity(action string) interface{} {
	return parser.Brief()
}

func (parser TransactionDetailCakeParser) DeleteActivity(action string) interface{} {
	return parser.Brief()
}

func (parser TransactionDetailCakeParser) GeneralActivity(action string) interface{} {
	return parser.Brief()
}

func (parser TransactionDetailCakeParser) UpdateActivity(action string) interface{} {
	return parser.Brief()
}

func (parser TransactionDetailCakeParser) Get() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, TransactionDetailCakeParser{Object: obj}.First())
	}
	return result
}

func (parser TransactionDetailCakeParser) First() interface{} {
	detailObj := parser.Object
	return map[string]interface{}{
		"id":            detailObj.ID,
		"transactionId": detailObj.TransactionID,
		"quantity":      detailObj.Quantity,
		"unitPrice":     detailObj.UnitPrice,
		"subTotal":      detailObj.SubTotal,
		"cake":          CakeParser{Object: detailObj.Cake}.Brief(),
		"createdAt":     detailObj.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":     detailObj.UpdatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser TransactionDetailCakeParser) Briefs() []interface{} {
	var result []interface{}
	for _, obj := range parser.Array {
		result = append(result, TransactionDetailCakeParser{Object: obj}.Brief())
	}
	return result
}

func (parser TransactionDetailCakeParser) Brief() interface{} {
	detailObj := parser.Object
	return map[string]interface{}{
		"id":        detailObj.ID,
		"cakeId":    detailObj.CakeID,
		"quantity":  detailObj.Quantity,
		"unitPrice": detailObj.UnitPrice,
		"subTotal":  detailObj.SubTotal,
		"createdAt": detailObj.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt": detailObj.UpdatedAt.Format("02/01/2006 15:04"),
	}
}
