package form

import (
	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
	"net/http"
	"service/internal/pkg/core"
)

type TestingForm struct {
	Name string   `json:"name"`
	Subs []string `json:"subs" validate:"required"`
}

func (rule *TestingForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(rule)
}

func (rule *TestingForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &rule)
}
