package web

import (
	"github.com/gorilla/mux"
	"service/internal/app/api/web/handler"
)

func Register(router *mux.Router) {
	activityRouter(router)
	testingRouter(router)
}

func activityRouter(router *mux.Router) {
	router.HandleFunc("/activities", handler.ActivityHandler{}.Get).Methods("GET")
}

func testingRouter(router *mux.Router) {
	var testingHandler handler.TestingHandler
	router.HandleFunc("/testings", testingHandler.Get).Methods("GET")
	router.HandleFunc("/testings", testingHandler.Create).Methods("POST")
	router.HandleFunc("/testings/upload/file", testingHandler.UploadByFile).Methods("POST")
	router.HandleFunc("/testings/upload/content", testingHandler.UploadByContent).Methods("POST")
}
