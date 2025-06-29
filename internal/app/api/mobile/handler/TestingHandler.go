package handler

import (
	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"net/http"
	"service/internal/pkg/parser"
	"service/internal/testing/repository"
)

type TestingHandler struct{}

func (ctr TestingHandler) Get(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewTestingRepository()
	testings, pagination, _ := repo.Find(r.URL.Query())

	psr := parser.TestingParser{Array: testings}

	res := xtremeres.Response{Array: psr.Get(), Pagination: &pagination}
	res.Success(w)
}
