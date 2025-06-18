package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

/*
ProductDashboard handles requests to the product dashboard.
Only users with "Admin", "User", or "SuperAdmin" roles for the product can access.
*/
func ProductDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	product := vars["product"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome to the " + product + " dashboard!",
	})
}
