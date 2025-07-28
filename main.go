package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Mryashbhardwaj/marketAnalysis/config"
	tradebook_service "github.com/Mryashbhardwaj/marketAnalysis/core/tradebook"
	"github.com/Mryashbhardwaj/marketAnalysis/routes"
)

func PrintCLIGreeting() {
	fmt.Println("====================================")
	fmt.Println("      Welcome to marketWatch 📈      ")
	fmt.Println("  Your CLI companion for tracking   ")
	fmt.Println("    mutual funds and equity data    ")
	fmt.Println("====================================")
	fmt.Println()
}

// initialise immemory database
func buildCache(cfg *config.Config) error {
	mfDir := cfg.MutualFunds.TradeFilesDirectory
	eqDir := cfg.Equity.TradeFilesDirectory

	// if neither is configured
	if mfDir == "" && eqDir == "" {
		return fmt.Errorf("no tradefiles directory set for equity or mutual funds")
	}

	if mfDir != "" {
		if err := tradebook_service.BuildMFTradeBook(mfDir); err != nil {
			return fmt.Errorf("MF TradeBook: %w", err)
		}
		if err := tradebook_service.BuildMFTrendCacheIfMissing(); err != nil {
			return fmt.Errorf("MF Trend cache: %w", err)
		}
		if err := tradebook_service.BuildMFPriceHistoryCacheFromFile(); err != nil {
			return fmt.Errorf("MF Price cache: %w", err)
		}

	}

	if eqDir != "" {
		if err := tradebook_service.BuildEquityTradeBook(eqDir); err != nil {
			return fmt.Errorf("EQ TradeBook: %w", err)
		}
		if err := tradebook_service.BuildEquityPriceHistoryCacheFromFile(); err != nil {
			return fmt.Errorf("EQ Price cache: %w", err)
		}
	}

	return nil
}

func main() {

	port := flag.Int64("p", 8080, "Port to run the server on")
	configFilePath := flag.String("c", "./config.yaml", "Location of the config file")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.LoadConfig(*configFilePath)
	if err != nil {
		logger.Error("failed to open config file", slog.String("error", err.Error()))
	}

	err = buildCache(cfg)
	if err != nil {
		logger.Error("failed building cache", slog.String("error", err.Error()))
		return
	}

	router := routes.SetupRouter()

	// initialise service
	host := ""
	addr := fmt.Sprintf("%s:%d", host, *port)
	logger.Info("Starting server ", slog.String("addr", addr))
	err = http.ListenAndServe(addr, router)
	if err != nil {
		fmt.Println("Failed to start HTTP server:", err)
	}
}
