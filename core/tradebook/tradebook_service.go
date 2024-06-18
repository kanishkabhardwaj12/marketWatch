package tradebook_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"time"

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

type MutualFundsTrade struct {
	Isin               string
	TradeDate          string
	Exchange           string
	Segment            string
	Series             string
	TradeType          string
	Auction            string
	Quantity           string
	Price              string
	TradeID            string
	OrderID            string
	OrderExecutionTime string
}

type FundName string
type ScriptName string

type MutualFundsTradebook struct {
	AllFunds             []FundName
	MutualFundsTradebook map[FundName][]MutualFundsTrade
}

type EquityTradebook struct {
	AllScripts      []ScriptName
	EquityTradebook map[ScriptName][]EquityTrade
}

var mutualFundsTradebook MutualFundsTradebook
var equityTradebook EquityTradebook

var shareHistory = make(map[ScriptName][]models.CandlePoint)

func readMFTradeFiles(tradebookDir string) (map[FundName][]MutualFundsTrade, error) {
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
	tradebook := make(map[FundName][]MutualFundsTrade)
	for _, record := range tradebookCSV {
		if record[0] == "symbol" {
			continue
		}
		if _, ok := tradeSet[record[10]]; ok {
			fmt.Println("found duplicate trade")
			continue
		}
		tradeSet[record[10]] = struct{}{}

		symbol := FundName(record[0])
		if _, ok := tradebook[symbol]; !ok {
			tradebook[symbol] = []MutualFundsTrade{}
		}
		tradebook[symbol] = append(tradebook[symbol], MutualFundsTrade{
			Isin:               record[1],
			TradeDate:          record[2],
			Exchange:           record[3],
			Segment:            record[4],
			Series:             record[5],
			TradeType:          record[6],
			Auction:            record[7],
			Quantity:           record[8],
			Price:              record[9],
			TradeID:            record[10],
			OrderID:            record[11],
			OrderExecutionTime: record[12],
		})

	}
	return tradebook, nil
}

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
			fmt.Println("found duplicate trade")
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

func BuildMFTradeBook(tradebookDir string) error {
	tradeMap, err := readMFTradeFiles(tradebookDir)
	if err != nil {
		return err
	}
	var trickers []FundName
	for fundName, _ := range tradeMap {
		trickers = append(trickers, fundName)
	}
	mutualFundsTradebook.MutualFundsTradebook = tradeMap
	mutualFundsTradebook.AllFunds = trickers
	return nil
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

func getFromMoneyControll(tickerSymbol ScriptName) (*models.MoneyControlResponse, error) {
	startTime := time.Unix(490147200, 0)
	endTime := time.Now()
	durationSince := math.Ceil(endTime.Sub(startTime).Hours() / 24)
	priceAPIURL := fmt.Sprintf("https://priceapi.moneycontrol.com/techCharts/indianMarket/stock/history?symbol=%s&resolution=1D&from=%d&to=%d&countback=%.f&currencyCode=INR", tickerSymbol, startTime.Unix(), endTime.Unix(), durationSince)
	fmt.Println(priceAPIURL)

	req, err := http.NewRequest("GET", priceAPIURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	k := models.MoneyControlResponse{}
	err = json.Unmarshal(body, &k)
	return &k, err
}

func persistInFile(symbol ScriptName, trend []models.CandlePoint) error {
	fileContent, err := json.Marshal(trend)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("./data/trends/equity/%s.json", symbol)
	return os.WriteFile(fileName, fileContent, os.ModePerm)
}

func fetchTradeHistories(script ScriptName) ([]models.CandlePoint, error) {
	k, err := getFromMoneyControll(script)
	if err != nil {
		return nil, err
	}

	candlePoints := make([]models.CandlePoint, len(k.T))

	for i, timeStamp := range k.T {
		candlePoints[i] = models.CandlePoint{
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
	for _, symbol := range equityTradebook.AllScripts {
		history, err := fetchTradeHistories(symbol)
		if err != nil {
			fmt.Printf("error fetching history for %s, err:%s", symbol, err.Error())
			continue
		}
		shareHistory[symbol] = history
		err = persistInFile(symbol, history)
		if err != nil {
			fmt.Printf("error persisting history for %s, err:%s", symbol, err.Error())
			continue
		}
	}
	return nil
}

func buildCacheFromFile(symbol ScriptName) ([]models.CandlePoint, error) {
	fileName := fmt.Sprintf("./data/trends/equity/%s.json", symbol)
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	trend := []models.CandlePoint{}
	err = json.Unmarshal(fileContent, &trend)
	return trend, err
}

func BuildPriceHistoryCacheFromFile() error {
	for _, symbol := range equityTradebook.AllScripts {
		history, err := buildCacheFromFile(symbol)
		if err != nil {
			fmt.Printf("error fetching history from file for %s, err:%s\n", symbol, err.Error())
			continue
		}
		shareHistory[symbol] = history
	}
	return nil
}

func GetMutualFundsList() []FundName {
	return mutualFundsTradebook.AllFunds
}

func GetEquityList() []ScriptName {
	return equityTradebook.AllScripts
}

func GetPriceTrend(symbol string) []models.CandlePoint {

	return shareHistory[ScriptName(symbol)]
}
