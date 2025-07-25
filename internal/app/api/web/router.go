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

	var ingredientHandler handler.IngredientHandler
	ingredientRouter := router.PathPrefix("/components/ingredients").Subrouter()
	ingredientRouter.HandleFunc("", ingredientHandler.Get).Methods("GET")
	ingredientRouter.HandleFunc("", ingredientHandler.Create).Methods("POST")
	ingredientRouter.HandleFunc("/{id}", ingredientHandler.Detail).Methods("GET")
	ingredientRouter.HandleFunc("/{id}", ingredientHandler.Update).Methods("PUT")
	ingredientRouter.HandleFunc("/{id}", ingredientHandler.Delete).Methods("DELETE")

	var cakeStaticHandler handler.CakeStaticHandler
	cakeStaticRouter := router.PathPrefix("/statics").Subrouter()
	cakeStaticRouter.HandleFunc("/units", cakeStaticHandler.GetUnitOfMeasure).Methods("GET")
	cakeStaticRouter.HandleFunc("/costs", cakeStaticHandler.GetCostType).Methods("GET")

	var cakeHandler handler.CakeHandler
	router.HandleFunc("", cakeHandler.Get).Methods("GET")
	router.HandleFunc("", cakeHandler.Create).Methods("POST")
	router.HandleFunc("/{id}", cakeHandler.Detail).Methods("GET")
	router.HandleFunc("/{id}", cakeHandler.Update).Methods("PUT")
	router.HandleFunc("/{id}", cakeHandler.Delete).Methods("DELETE")

}

func transactionRouter(router *mux.Router) {
	router = router.PathPrefix("/transactions").Subrouter()

	var transactionHandler handler.TransactionHandler
	router.HandleFunc("", transactionHandler.Get).Methods("GET")
	router.HandleFunc("", transactionHandler.Create).Methods("POST")
	router.HandleFunc("/{id}", transactionHandler.Detail).Methods("GET")
	router.HandleFunc("/{id}", transactionHandler.Update).Methods("PUT")
	router.HandleFunc("/{id}", transactionHandler.Delete).Methods("DELETE")
	router.HandleFunc("/download/excel", transactionHandler.ReportExcel).Methods("POST")
}
