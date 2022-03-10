package routers

import (
	"go-postgres/middlewares"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/user/{id}", middlewares.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newuser", middlewares.CreateUser).Methods("POST", "OPTIONS")

	return router
}
