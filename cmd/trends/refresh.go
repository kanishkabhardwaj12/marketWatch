package trends

import (
	"errors"
	"log/slog"
	"os"

	"github.com/Mryashbhardwaj/marketAnalysis/internal/config"
	"github.com/Mryashbhardwaj/marketAnalysis/internal/domain/service"
	"github.com/spf13/cobra"
)

type fetchCommand struct {
	configFilePath string

	logger *slog.Logger
	config *config.Config
}

// NewFetchCommand initializes command to fetch the API
func NewFetchCommand() *cobra.Command {
	s := &fetchCommand{}

	cmd := &cobra.Command{
		Use:     "refresh-trends",
		Short:   "refresh trends from the API",
		Example: "marketWatch fetch-trends -c /path/to/config.yaml",
		RunE:    s.RunE,
		PreRunE: s.PreRunE,
	}

	// Config filepath flag
	cmd.Flags().StringVarP(&s.configFilePath, "config", "c", "", "File path for client configuration")

	return cmd
}

func (s *fetchCommand) PreRunE(cmd *cobra.Command, _ []string) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	s.logger = logger

	if s.configFilePath == "" {
		return errors.New("config file path is required")
	}

	cfg, err := config.LoadConfig(s.configFilePath)
	if err != nil {
		logger.Error("failed to open config file", slog.String("error", err.Error()))
		return err
	}

	s.config = cfg

	return err
}

func (s *fetchCommand) RunE(_ *cobra.Command, _ []string) error {

	err := service.BuildCache(s.config)
	if err != nil {
		s.logger.Error("failed building cache", slog.String("error", err.Error()))
		return err
	}
	err = service.BuildPriceHistoryCache()
	if err != nil {
		s.logger.Error("failed building price history cache for Equity", slog.String("error", err.Error()))
		return err
	}
	err = service.BuildMFPriceHistoryCache()
	if err != nil {
		s.logger.Error("failed building price history cache for Mutual Funds", slog.String("error", err.Error()))
		return err
	}

	return nil
}
