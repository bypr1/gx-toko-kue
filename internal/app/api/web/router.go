package web

import (
	"service/internal/app/api/web/handler"

	"github.com/gorilla/mux"
)

func Register(router *mux.Router) {
	activityRouter(router)
	testingRouter(router)
	cakeRouter(router)
	ingredientRouter(router)
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
	var cakeHandler handler.CakeHandler
	router.HandleFunc("/cakes", cakeHandler.Get).Methods("GET")
	router.HandleFunc("/cakes", cakeHandler.Create).Methods("POST")
	router.HandleFunc("/cakes/{id}", cakeHandler.Detail).Methods("GET")
	router.HandleFunc("/cakes/{id}", cakeHandler.Update).Methods("PUT")
	router.HandleFunc("/cakes/{id}", cakeHandler.Delete).Methods("DELETE")
	router.HandleFunc("/cakes/{id}/calculate-cost", cakeHandler.CalculateCost).Methods("GET")
}

func ingredientRouter(router *mux.Router) {
	var ingredientHandler handler.IngredientHandler
	router.HandleFunc("/ingredients", ingredientHandler.Get).Methods("GET")
	router.HandleFunc("/ingredients", ingredientHandler.Create).Methods("POST")
	router.HandleFunc("/ingredients/{id}", ingredientHandler.Detail).Methods("GET")
	router.HandleFunc("/ingredients/{id}", ingredientHandler.Update).Methods("PUT")
	router.HandleFunc("/ingredients/{id}", ingredientHandler.Delete).Methods("DELETE")
}
