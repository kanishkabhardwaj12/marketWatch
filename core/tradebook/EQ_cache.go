package tradebook_service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	MC "github.com/Mryashbhardwaj/marketAnalysis/clients/moneyControl"
	"github.com/Mryashbhardwaj/marketAnalysis/models"
	"github.com/Mryashbhardwaj/marketAnalysis/utils"
)

type EquityTrade struct {
	Isin               string
	Symbol             string
	TradeDate          string
	Exchange           string
	Segment            string
	Series             string
	TradeType          string
	Auction            string
	Quantity           string
	Price              string
	TradeId            string
	OrderId            string
	OrderExecutionTime string
}

type ScriptName string

func (s ScriptName) String() string {
	return string(s)
}

type EquityTradebook struct {
	AllScripts      []ScriptName
	EquityTradebook map[ScriptName][]EquityTrade
}

var equityTradebook EquityTradebook

var shareHistory = make(map[ScriptName][]models.EquityPriceData)

func readEquityTradeFiles(tradebookDir string) (map[ScriptName][]EquityTrade, error) {
	// to remove duplidate trade ids
	tradeSet := make(map[string]struct{})

	tradeFiles, err := utils.ReadDir(tradebookDir)
	if err != nil {
		return nil, err
	}
	tradebookCSV, err := utils.ReadCSV(tradeFiles)
	if err != nil {
		return nil, err
	}
	tradebook := make(map[ScriptName][]EquityTrade)
	for _, record := range tradebookCSV {
		if record[1] == "symbol" {
			continue
		}
		if _, ok := tradeSet[record[10]]; ok {
			continue
		}
		tradeSet[record[10]] = struct{}{}

		symbol := ScriptName(record[0])
		if _, ok := tradebook[symbol]; !ok {
			tradebook[symbol] = []EquityTrade{}
		}
		tradebook[symbol] = append(tradebook[symbol], EquityTrade{
			Symbol:             record[1],
			TradeDate:          record[2],
			Exchange:           record[3],
			Segment:            record[4],
			TradeType:          record[6],
			Quantity:           record[8],
			Price:              record[9],
			OrderExecutionTime: record[12],
		})
	}
	return tradebook, nil
}

func BuildEquityTradeBook(tradebookDir string) error {
	tradeMap, err := readEquityTradeFiles(tradebookDir)
	if err != nil {
		return err
	}
	var trickers []ScriptName
	for fundName, _ := range tradeMap {
		trickers = append(trickers, fundName)
	}
	equityTradebook.EquityTradebook = tradeMap
	equityTradebook.AllScripts = trickers
	return nil
}

func persistInFile(symbol string, trend interface{}) error {
	fileContent, err := json.Marshal(trend)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("./data/trends/EQ/%s.json", symbol)
	return os.WriteFile(fileName, fileContent, os.ModePerm)
}

func fetchTradeHistories(script ScriptName) ([]models.EquityPriceData, error) {
	k, err := MC.GetEQHistoryFromMoneyControll(script.String())
	if err != nil {
		return nil, err
	}

	candlePoints := make([]models.EquityPriceData, len(k.T))

	for i, timeStamp := range k.T {
		candlePoints[i] = models.EquityPriceData{
			Close:      k.C[i],
			High:       k.H[i],
			Volume:     k.V[i],
			Open:       k.O[i],
			Low:        k.L[i],
			Timestamps: time.Unix(timeStamp, 0),
		}
	}
	return candlePoints, nil
}

func BuildPriceHistoryCache() error {
<<<<<<< Updated upstream:core/tradebook/EQ_cache.go
	var errorList []string
	for _, symbol := range equityTradebook.AllScripts {
		history, err := fetchTradeHistories(symbol)
		if err != nil {
			fmt.Printf("error fetching history for %s, err:%s", symbol, err.Error())
			errorList = append(errorList, fmt.Sprintf("error fetching history for %s, err:%s", symbol, err.Error()))
			continue
		}
		shareHistory[symbol] = history
		err = persistInFile(string(symbol), history)
		if err != nil {
			fmt.Printf("error persisting history for %s, err:%s", symbol, err.Error())
			errorList = append(errorList, fmt.Sprintf("error persisting history for %s, err:%s", symbol, err.Error()))
			continue
		}
=======
	var (
		mu        sync.Mutex
		errorList []string
		wg        sync.WaitGroup
	)

	semaphore := make(chan struct{}, 5) //5 concurrent requests

	for _, symbol := range EquityTradebookCache.AllScripts {
		wg.Add(1)

		go func(sym ScriptName) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			fmt.Printf("Processing EQ: %s\n", sym)

			history, err := fetchTradeHistories(sym)
			if err != nil {
				mu.Lock()
				errorList = append(errorList, fmt.Sprintf("error fetching history for %s, err:%s", sym, err.Error()))
				mu.Unlock()
				return
			}
			shareHistory[sym] = history

			err = persistInFile(string(sym), history)
			if err != nil {
				mu.Lock()
				errorList = append(errorList, fmt.Sprintf("error persisting history for %s, err:%s", sym, err.Error()))
				mu.Unlock()
			}
		}(symbol)
>>>>>>> Stashed changes:internal/domain/service/EQ_cache.go
	}

	wg.Wait()

	if len(errorList) == 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(errorList, "\n"))
}

func buildEquityCacheFromFile(symbol ScriptName) ([]models.EquityPriceData, error) {
	fileName := fmt.Sprintf("./data/trends/EQ/%s.json", symbol)
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	trend := []models.EquityPriceData{}
	err = json.Unmarshal(fileContent, &trend)
	return trend, err
}

func BuildEquityPriceHistoryCacheFromFile() error {
	for _, symbol := range equityTradebook.AllScripts {
		history, err := buildEquityCacheFromFile(symbol)
		if err != nil {
			fmt.Printf("error fetching history from file for %s, err:%s\n", symbol, err.Error())
			continue
		}
		shareHistory[symbol] = history
	}
	return nil
}
