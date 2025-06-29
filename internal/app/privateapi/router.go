package privateapi

import (
	"github.com/gorilla/mux"
	"service/internal/app/privateapi/handler"
	"service/internal/pkg/middleware"
)

func Register(router *mux.Router) {
	api := router.PathPrefix("/private-api").Subrouter()
	api.Use(middleware.AuthenticatePrivateAPI)

	testingRouter(api)
}

func testingRouter(router *mux.Router) {
	var testingHandler handler.TestingHandler
	router.HandleFunc("/testings", testingHandler.Get).Methods("GET")
	router.HandleFunc("/testings", testingHandler.Create).Methods("POST")
	router.HandleFunc("/testings/upload/file", testingHandler.UploadByFile).Methods("POST")
	router.HandleFunc("/testings/upload/content", testingHandler.UploadByContent).Methods("POST")
}
