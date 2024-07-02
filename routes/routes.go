package routes

import (
	"github.com/Mryashbhardwaj/marketAnalysis/handlers"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/equity/list", handlers.GetEquityList).Methods("GET")
	router.HandleFunc("/api/trend", handlers.GetTrend).Methods("GET")
	router.HandleFunc("/api/trend/compare", handlers.GetTrendComparison).Methods("GET")
	router.HandleFunc("/api/equity/history/refresh", handlers.RefreshPriceHistory).Methods("GET")

	router.HandleFunc("/api/mf/list", handlers.GetMutualFundsList).Methods("GET")
	router.HandleFunc("/api/mf/positions", handlers.GetMFPositions).Methods("GET")
	router.HandleFunc("/api/mf/trend", handlers.GetMFTrend).Methods("GET")
	router.HandleFunc("/api/mf/summary", handlers.GetMFSummary).Methods("GET")
	router.HandleFunc("/api/mf/trend/compare", handlers.GetMFGrowthComparison).Methods("GET")
	router.HandleFunc("/api/mf/history/refresh", handlers.RefreshMFPriceHistory).Methods("GET")

	return router
}

// important links
//  https://console.zerodha.com/reports/tradebook
//  link to refresh the trends data
