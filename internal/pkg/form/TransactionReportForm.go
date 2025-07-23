package form

import (
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type TransactionReportForm struct {
	FromDate  string   `json:"fromDate" validate:"required"`
	ToDate    string   `json:"toDate" validate:"required"`
	MinAmount *float64 `json:"minAmount" validate:"omitempty,gt=0"`
	MaxAmount *float64 `json:"maxAmount" validate:"omitempty,gt=0"`
}

func (f *TransactionReportForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(f)
}

func (f *TransactionReportForm) APIParse(r *http.Request) {
	form := core.BaseForm{}
	form.APIParse(r, &f)
}
