package form

import (
	"net/http"
	"service/internal/pkg/core"
	"time"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type TransactionForm struct {
	TransactionDate string                `json:"transactionDate" validate:"required"`
	Cakes           []TransactionCakeForm `json:"cakes" validate:"required,dive"`
}

type TransactionCakeForm struct {
	CakeID   uint `json:"cakeId" validate:"required,gt=0"`
	Quantity int  `json:"quantity" validate:"required,gt=0"`
}

func (f *TransactionForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(f)
}

func (f *TransactionForm) APIParse(r *http.Request) {
	form := core.BaseForm{}
	form.APIParse(r, &f)
}

func (f *TransactionForm) GetTransactionDate() time.Time {
	date, err := time.Parse("2006-01-02", f.TransactionDate)
	if err != nil {
		return time.Now()
	}
	return date
}
