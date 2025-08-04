package service

import (
	"fmt"

	"github.com/Mryashbhardwaj/marketAnalysis/internal/config"
)

// initialise immemory database
func BuildCache(cfg *config.Config) error { // move this function to utils package
	mfDir := cfg.MutualFunds.TradeFilesDirectory
	eqDir := cfg.Equity.TradeFilesDirectory

	// if neither is configured
	if mfDir == "" && eqDir == "" {
		return fmt.Errorf("no tradefiles directory set for equity or mutual funds")
	}

	if mfDir != "" {
		if err := BuildMFTradeBook(mfDir); err != nil {
			return fmt.Errorf("MF TradeBook: %w", err)
		}
		if err := BuildMFTrendCacheIfMissing(); err != nil {
			return fmt.Errorf("MF Trend cache: %w", err)
		}
		if err := BuildMFPriceHistoryCacheFromFile(); err != nil {
			return fmt.Errorf("MF Price cache: %w", err)
		}

	}

	if eqDir != "" {
		if err := BuildEquityTradeBook(eqDir); err != nil {
			return fmt.Errorf("EQ TradeBook: %w", err)
		}
		if err := BuildEquityPriceHistoryCacheFromFile(); err != nil {
			return fmt.Errorf("EQ Price cache: %w", err)
		}
	}

	return nil
}
