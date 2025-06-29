package handler

import (
	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"net/http"
	repository2 "service/internal/activity/repository"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/parser"
	"service/internal/testing/service"
)

type TestingHandler struct{}

func (ctr TestingHandler) Get(w http.ResponseWriter, r *http.Request) {
	//repo := repository.NewTestingRepository()
	//testings, pagination, _ := repo.Find(r.URL.Query())
	//
	//psr := parser.TestingParser{Array: testings}
	var testings []interface{}
	testings = append(testings, map[string]interface{}{"testing": "test"})

	res := xtremeres.Response{Array: testings}
	res.Success(w)
}

func (ctr TestingHandler) Create(w http.ResponseWriter, r *http.Request) {
	form := form2.TestingForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewTestingService()
	srv.SetActivityRepository(repository2.NewActivityRepository())

	testing := srv.Create(form)

	psr := parser.TestingParser{Object: testing}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (ctr TestingHandler) UploadByFile(w http.ResponseWriter, r *http.Request) {
	form := form2.TestingUploadForm{}
	form.APIParse(r)

	srv := service.NewTestingService()
	uploaded := srv.UploadByFile(form)

	res := xtremeres.Response{Object: uploaded}
	res.Success(w)
}

func (ctr TestingHandler) UploadByContent(w http.ResponseWriter, r *http.Request) {
	form := form2.TestingUploadContentForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewTestingService()
	uploaded := srv.UploadByContent(form)

	res := xtremeres.Response{Object: uploaded}
	res.Success(w)
}
