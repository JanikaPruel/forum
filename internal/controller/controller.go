package controller

import "net/http"

// MainContrller
func MainContrller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("main controller"))
}
