package routes

import (
	"github.com/Mryashbhardwaj/marketAnalysis/handlers"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/trend", handlers.GetTrend).Methods("GET")
	router.HandleFunc("/api/mutual_funds", handlers.GetMutualFundsList).Methods("GET")
	// router.HandleFunc("/api/users", handlers.CreateUser).Methods("POST")
	// router.HandleFunc("/api/users/{id}", handlers.UpdateUser).Methods("PUT")
	// router.HandleFunc("/api/users/{id}", handlers.DeleteUser).Methods("DELETE")

	return router
}
