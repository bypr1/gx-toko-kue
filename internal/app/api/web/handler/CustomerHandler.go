package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	f "service/internal/pkg/form"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

type CustomerHandler struct{}

func (CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var form f.CustomerForm
	form.APIParse(r)
	form.Validate()

	// Convert struct to JSON
	jsonData, err := json.Marshal(form)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return
	}

	res := xtremeres.Response{Object: result}
	res.Success(w)
}
