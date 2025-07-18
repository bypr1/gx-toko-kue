package error

import (
	"net/http"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

func ErrXtremeTransactionGet(internalMsg string) {
	xtremeres.Error(http.StatusNotFound, "Transaction not found", internalMsg, false, nil)
}

func ErrXtremeTransactionSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save transaction", internalMsg, false, nil)
}

func ErrXtremeTransactionDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete transaction", internalMsg, false, nil)
}

func ErrXtremeTransactionCakeSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save transaction detail", internalMsg, false, nil)
}

func ErrXtremeTransactionCakeGet(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to retrieve transactions", internalMsg, false, nil)
}

func ErrXtremeTransactionCakeDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete transaction detail", internalMsg, false, nil)
}

func ErrXtremeTransactionExcelGenerate(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to generate transaction Excel file", internalMsg, false, nil)
}
