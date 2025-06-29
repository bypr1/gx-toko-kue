package api

import (
	"fmt"
	"net/http"
	"os"
	"service/internal/app/api/mobile"
	"service/internal/app/api/web"

	"github.com/gorilla/mux"
)

func Register(router *mux.Router) {
	version := os.Getenv("VERSION")
	service := os.Getenv("SERVICE")

	api := router.PathPrefix("/api").Subrouter()
	//api.Use(middleware.EmployeeIdentifier) // TODO: Re-enable this code after installing github.com/globalxtreme/go-identifier module (If you use GX Identifier for authorization)

	web.Register(api.PathPrefix(fmt.Sprintf("/web/%s/%s", version, service)).Subrouter())
	mobile.Register(api.PathPrefix(fmt.Sprintf("/mobile/%s/%s", version, service)).Subrouter())

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintf(w, "Welcome to %s API version %s", service, version); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
