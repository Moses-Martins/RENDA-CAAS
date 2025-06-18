package main

import (
	"RENDA-CAAS/config"
	"RENDA-CAAS/controllers"
	"RENDA-CAAS/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()
	controllers.InitUserCollection()

	r := mux.NewRouter()
	r.HandleFunc("/v1/me", controllers.Me).Methods("GET")
	r.HandleFunc("/v1/register/renda360", controllers.RegisterRenda360).Methods("POST")
	r.HandleFunc("/v1/register/scale", controllers.RegisterScale).Methods("POST")
	r.HandleFunc("/v1/register/horizon", controllers.RegisterHorizon).Methods("POST")
	r.HandleFunc("/v1/login", controllers.Login).Methods("POST")
	r.HandleFunc("/v1/admin/update-privilege", controllers.UpdateUserPrivilege).Methods("PATCH")
	r.Handle("/v1/dashboard/{product}", middleware.AdminOrUserForProduct(controllers.UserCollection)(http.HandlerFunc(controllers.ProductDashboard))).Methods("GET")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
