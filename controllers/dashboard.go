package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func ProductDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	product := vars["product"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome to the " + product + " dashboard!",
	})
}
