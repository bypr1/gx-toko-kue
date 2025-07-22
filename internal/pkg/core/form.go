package core

import (
	"encoding/json"
	"net/http"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
	formParser "github.com/go-playground/form/v4"
)

type FormInterface interface {
	Validate()
}

type APIFormInterface interface {
	APIParse(r *http.Request)
}

type BaseForm struct{}

func (BaseForm) APIParse(r *http.Request, rule interface{}) interface{} {
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		xtremeres.ErrXtremeBadRequest(err.Error())
	}

	return rule
}

func (BaseForm) FormParse(r *http.Request, form interface{}) interface{} {
	decoder := formParser.NewDecoder()

	formValue := r.MultipartForm.Value
	convertedFormValue := make(map[string][]string, 0)
	for key, value := range formValue {
		var convertedKey string
		lastIndex := 0
		isNumber := false
		for i := 0; i < len(key); i++ {
			switch key[i] {
			case '[':
				convertedKey += key[lastIndex:i]
				lastIndex = i + 1
				isNumber = true
			case ']':
				if !isNumber {
					convertedKey += "."
					convertedKey += key[lastIndex:i]
				} else {
					convertedKey += "["
					convertedKey += key[lastIndex:i]
					convertedKey += "]"
				}
				lastIndex = i + 1
			default:
				isNumber = isNumber && key[i] >= '0' && key[i] <= '9'
			}
		}

		if convertedKey == "" {
			convertedKey = key
		}

		convertedFormValue[convertedKey] = value
	}

	if err := decoder.Decode(&form, convertedFormValue); err != nil {
		xtremeres.ErrXtremeBadRequest(err.Error())
	}

	return form
}
