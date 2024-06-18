package routes

import (
	"github.com/Mryashbhardwaj/marketAnalysis/handlers"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/trend", handlers.GetTrend).Methods("GET")
	router.HandleFunc("/api/mutual_funds/list", handlers.GetMutualFundsList).Methods("GET")
	router.HandleFunc("/api/equity/list", handlers.GetEquityList).Methods("GET")
	router.HandleFunc("/api/symbols/history/refresh", handlers.RefreshPriceHistory).Methods("POST")

	return router
}
