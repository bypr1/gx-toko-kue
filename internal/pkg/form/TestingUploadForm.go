package form

import (
	"net/http"
)

type TestingUploadForm struct {
	Request *http.Request
}

func (rule *TestingUploadForm) APIParse(r *http.Request) {
	rule.Request = r
}
