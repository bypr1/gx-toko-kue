package web

import (
	"service/internal/app/api/web/handler"

	"github.com/gorilla/mux"
)

func Register(router *mux.Router) {
	activityRouter(router)
	testingRouter(router)
	cakeRouter(router)
	transactionRouter(router)
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

func cakeRouter(router *mux.Router) {
	router = router.PathPrefix("/cakes").Subrouter()

	var cakeHandler handler.CakeHandler
	router.HandleFunc("", cakeHandler.Get).Methods("GET")
	router.HandleFunc("", cakeHandler.Create).Methods("POST")
	router.HandleFunc("/{id}", cakeHandler.Detail).Methods("GET")
	router.HandleFunc("/{id}", cakeHandler.Update).Methods("PUT")
	router.HandleFunc("/{id}", cakeHandler.Delete).Methods("DELETE")

	var ingredientHandler handler.IngredientHandler
	router.HandleFunc("/components/ingredients", ingredientHandler.Get).Methods("GET")
	router.HandleFunc("/components/ingredients", ingredientHandler.Create).Methods("POST")
	router.HandleFunc("/components/ingredients/{id}", ingredientHandler.Detail).Methods("GET")
	router.HandleFunc("/components/ingredients/{id}", ingredientHandler.Update).Methods("PUT")
	router.HandleFunc("/components/ingredients/{id}", ingredientHandler.Delete).Methods("DELETE")

}

func transactionRouter(router *mux.Router) {
	router = router.PathPrefix("/transactions").Subrouter()

	var transactionHandler handler.TransactionHandler
	router.HandleFunc("", transactionHandler.Get).Methods("GET")
	router.HandleFunc("", transactionHandler.Create).Methods("POST")
	router.HandleFunc("/{id}", transactionHandler.Detail).Methods("GET")
	router.HandleFunc("/{id}", transactionHandler.Update).Methods("PUT")
	router.HandleFunc("/{id}", transactionHandler.Delete).Methods("DELETE")
}
