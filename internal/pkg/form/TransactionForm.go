package form

import (
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type TransactionForm struct {
	TransactionDate string                `json:"transactionDate" validate:"required"`
	Cakes           []TransactionCakeForm `json:"cakes" validate:"required,dive"`
}

type TransactionCakeForm struct {
	ID       uint `json:"id"`
	CakeID   uint `json:"cakeId" validate:"required,gt=0"`
	Quantity int  `json:"quantity" validate:"required,gt=0"`
	Deleted  bool `json:"deleted"`
}

func (f *TransactionForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(f)
}

func (f *TransactionForm) APIParse(r *http.Request) {
	form := core.BaseForm{}
	form.APIParse(r, &f)
}
