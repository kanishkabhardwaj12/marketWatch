package routes

import (
	"github.com/Mryashbhardwaj/marketAnalysis/internal/api/handlers"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// generate a single endpoint for mf and eq
	// list, comparison, trend

	router.HandleFunc("/api/equity/list", handlers.GetEquityList).Methods("GET")
	router.HandleFunc("/api/equity/trend", handlers.GetTrend).Methods("GET")
	router.HandleFunc("/api/equity/trend/compare", handlers.GetTrendComparison).Methods("GET")
	router.HandleFunc("/api/equity/history/refresh", handlers.RefreshPriceHistory).Methods("GET")
	router.HandleFunc("/api/equity/breakdown", handlers.GetEqBreakdown).Methods("GET")

	router.HandleFunc("/api/mutual_funds/list", handlers.GetMutualFundsList).Methods("GET")
	router.HandleFunc("/api/mutual_funds/positions", handlers.GetMFPositions).Methods("GET")
	router.HandleFunc("/api/mutual_funds/trend", handlers.GetMFTrend).Methods("GET")
	router.HandleFunc("/api/mutual_funds/summary", handlers.GetMFSummary).Methods("GET")
	router.HandleFunc("/api/mutual_funds/trend/compare", handlers.GetMFGrowthComparison).Methods("GET")
	router.HandleFunc("/api/mutual_funds/history/refresh", handlers.RefreshMFPriceHistory).Methods("GET")

	return router
}

// important links
//  https://console.zerodha.com/reports/tradebook
//  link to refresh the trends data
