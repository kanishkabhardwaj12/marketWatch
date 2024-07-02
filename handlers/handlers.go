package handlers

import (
	"net/http"
	"strings"
	"time"

	tradebook_service "github.com/Mryashbhardwaj/marketAnalysis/core/tradebook"
	"github.com/Mryashbhardwaj/marketAnalysis/utils"
)

type mutualFundOverview struct {
	Symbol        string
	Name          string
	CurrentTotal  float64
	HoldingsStart time.Time
	HoldingsEnd   time.Time
}

func GetTrend(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	from, to, err := utils.GetTimeRange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.RespondWithJSON(w, 200, tradebook_service.GetPriceTrendInTimeRange(symbol, from, to))
}

func GetMFSummary(w http.ResponseWriter, r *http.Request) {
	from, to, err := utils.GetTimeRange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.RespondWithJSON(w, 200, tradebook_service.GetMFSummmary(from, to))
}

func GetMFTrend(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	from, to, err := utils.GetTimeRange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.RespondWithJSON(w, 200, tradebook_service.GetPriceMFTrendInTimeRange(symbol, from, to))
}

func GetMFPositions(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	from, to, err := utils.GetTimeRange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.RespondWithJSON(w, 200, tradebook_service.GetPriceMFPositionsInTimeRange(symbol, from, to))
}

func GetTrendComparison(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	//TODO:  add condition here that if symbol is single value then give empty result
	symbols := strings.Split(symbol[1:len(symbol)-1], ",")

	from, to, err := utils.GetTimeRange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.RespondWithJSON(w, 200, tradebook_service.GetGrowthComparison(symbols, from, to))
}

func GetMFGrowthComparison(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	symbols := strings.Split(symbol[1:len(symbol)-1], ",")

	from, to, err := utils.GetTimeRange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.RespondWithJSON(w, 200, tradebook_service.GetMFGrowthComparison(symbols, from, to))
}

func GetMutualFundsList(w http.ResponseWriter, r *http.Request) {
	mfList := tradebook_service.GetMutualFundsList()
	utils.RespondWithJSON(w, 200, mfList)
}

func GetEquityList(w http.ResponseWriter, r *http.Request) {
	mfList := tradebook_service.GetEquityList()
	utils.RespondWithJSON(w, 200, mfList)
}

func RefreshPriceHistory(w http.ResponseWriter, r *http.Request) {
	err := tradebook_service.BuildPriceHistoryCache()
	if err != nil {
		utils.RespondWithJSON(w, 505, err.Error())
		return
	}
	utils.RespondWithJSON(w, 200, "Price History Refreshed Successfully")
}

func RefreshMFPriceHistory(w http.ResponseWriter, r *http.Request) {
	err := tradebook_service.BuildMFPriceHistoryCache()
	if err != nil {
		utils.RespondWithJSON(w, 505, err.Error())
		return
	}
	utils.RespondWithJSON(w, 200, "Price History Refreshed Successfully")
}
