package main

import (
	"fmt"
	"log"
	"net/http"

	tradebook_service "github.com/Mryashbhardwaj/marketAnalysis/core/tradebook"
	"github.com/Mryashbhardwaj/marketAnalysis/routes"
)

func buildCache() error {
	err := tradebook_service.BuildMFTradeBook("./data/trade_books/MF/")
	if err != nil {
		return err
	}
	err = tradebook_service.BuildEquityTradeBook("./data/trade_books/EQ/")
	if err != nil {
		return err
	}

	err = tradebook_service.BuildPriceHistoryCacheFromFile()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := buildCache()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	router := routes.SetupRouter()

	// initialise service
	log.Println("Starting server on :80")
	err = http.ListenAndServe(":80", router)
	if err != nil {
		fmt.Println("Failed to start HTTP server:", err)
	}
}
