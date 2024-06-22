package handlers

import (
	"fmt"
	"net/http"
	"strconv"
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
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	fmt.Println(fromStr, toStr)
	if fromStr == "" {
		fromStr = "490147200000"
	}
	if toStr == "" {
		toStr = strconv.FormatInt(time.Now().UnixMilli(), 10)
	}

	fmt.Println(fromStr, toStr)
	// Parse the "from" timestamp
	fromMilli, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid 'from' timestamp format. Expected epoch milliseconds", http.StatusBadRequest)
		return
	}
	from := time.Unix(fromMilli/1000, 0)

	// Parse the "to" timestamp
	toMilli, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid 'to' timestamp format. Expected epoch milliseconds", http.StatusBadRequest)
		return
	}
	to := time.Unix(toMilli/1000, 0)

	fmt.Println(from, to)
	utils.RespondWithJSON(w, 200, tradebook_service.GetPriceTrendInTimeRange(symbol, from, to))
}

func GetTrendComparison(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")

	symbols := strings.Split(symbol[1:len(symbol)-1], ",")

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	fmt.Println(fromStr, toStr)
	if fromStr == "" {
		fromStr = "490147200000"
	}
	if toStr == "" {
		toStr = strconv.FormatInt(time.Now().UnixMilli(), 10)
	}

	fmt.Println(fromStr, toStr)
	// Parse the "from" timestamp
	fromMilli, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid 'from' timestamp format. Expected epoch milliseconds", http.StatusBadRequest)
		return
	}
	from := time.Unix(fromMilli/1000, 0)

	// Parse the "to" timestamp
	toMilli, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid 'to' timestamp format. Expected epoch milliseconds", http.StatusBadRequest)
		return
	}
	to := time.Unix(toMilli/1000, 0)

	fmt.Println(from, to)
	utils.RespondWithJSON(w, 200, tradebook_service.GetGrowthComparison(symbols, from, to))
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
	// fetch share histories
	err := tradebook_service.BuildPriceHistoryCache()
	// Refactor to save price history in files
	if err != nil {
		utils.RespondWithJSON(w, 505, err.Error())
		return
	}
	utils.RespondWithJSON(w, 200, "Price History Refreshed Successfully")
}
