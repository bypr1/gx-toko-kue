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

func ErrXtremeTransactionDetailSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save transaction detail", internalMsg, false, nil)
}

func ErrXtremeTransactionDetailDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete transaction detail", internalMsg, false, nil)
}
