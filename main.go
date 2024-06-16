package main

import (
	"fmt"
	"log"
	"net/http"

	tradebook_service "github.com/Mryashbhardwaj/marketAnalysis/core/tradebook"
	"github.com/Mryashbhardwaj/marketAnalysis/routes"
)

func buildCache() error {
	return tradebook_service.BuildTradeBook("./data/trade_books/mutual_funds/")
}

func main() {
	err := buildCache()
	if err != nil {
		fmt.Println(err.Error())
	}

	router := routes.SetupRouter()

	// initialise service
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}
