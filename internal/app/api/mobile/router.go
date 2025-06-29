package mobile

import (
	"github.com/gorilla/mux"
	"service/internal/app/api/mobile/handler"
)

func Register(router *mux.Router) {

	// Testing
	var testingHandler handler.TestingHandler
	router.HandleFunc("/testings", testingHandler.Get).Methods("GET")
}
